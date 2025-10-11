package serverd

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-common/options"
	"github.com/xmx/aegis-common/tunnel/tundial"
	"github.com/xmx/aegis-common/tunnel/tunutil"
	"github.com/xmx/aegis-control/contract/authmesg"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"github.com/xmx/aegis-server/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func New(repo repository.All, cfg *config.Config, opts ...options.Lister[option]) tunutil.Handler {
	opt := options.Eval[option](opts...)
	if opt.huber == nil {
		opt.huber = linkhub.NewHub(16)
	}
	if opt.parent == nil {
		opt.parent = context.Background()
	}

	return &brokerServer{
		repo: repo,
		cfg:  cfg,
		opt:  opt,
	}
}

type brokerServer struct {
	repo repository.All
	cfg  *config.Config
	opt  option
}

func (as *brokerServer) Handle(mux tundial.Muxer) {
	defer mux.Close()

	if !as.allow() {
		as.log().Warn("限流器抑制 broker 上线")
		return
	}

	const defaultTimeout = 30 * time.Second
	_, peer, succeed := as.authentication(mux, defaultTimeout)
	if !succeed {
		return
	}
	defer as.disconnected(peer, defaultTimeout)

	srv := as.getServer(peer)
	_ = srv.Serve(mux)
}

// authentication 节点认证。
// 客户端主动建立一条虚拟子流连接用于交换认证信息，认证后改子流关闭，后续子流即为业务流。
func (as *brokerServer) authentication(mux tundial.Muxer, timeout time.Duration) (*authmesg.BrokerToServerRequest, linkhub.Peer, bool) {
	protocol, subprotocol := mux.Protocol()
	localAddr := mux.Addr().String()
	remoteAddr := mux.RemoteAddr().String()
	attrs := []any{
		slog.String("protocol", protocol), slog.String("subprotocol", subprotocol),
		slog.String("local_addr", localAddr), slog.String("remote_addr", remoteAddr),
	}

	// 设置超时主动断开，防止恶意客户端一直不建立认证连接。
	timer := time.AfterFunc(timeout, func() { _ = mux.Close() })
	defer timer.Stop()

	sig, err := mux.Accept()
	timer.Stop()
	if err != nil {
		as.log().Error("等待客户端建立认证连接出错", "error", err)
		return nil, nil, false
	}
	defer sig.Close()

	// 读取数据
	now := time.Now()
	_ = sig.SetDeadline(now.Add(timeout))
	req := new(authmesg.BrokerToServerRequest)
	if err = tunutil.ReadHead(sig, req); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		as.log().Error("读取请求信息错误", attrs...)
		return nil, nil, false
	}

	attrs = append(attrs, slog.Any("broker_auth_request", req))
	if err = as.validate(req); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		as.log().Error("读取请求信息校验错误", attrs...)
		_ = as.writeError(sig, http.StatusBadRequest, err)
		return nil, nil, false
	}
	broker, err := as.checkout(req, timeout)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
		as.log().Error("查询 broker 节点错误", attrs...)
		_ = as.writeError(sig, http.StatusInternalServerError, err)
		return nil, nil, false
	}
	if broker.Status { // 节点已经在线了
		as.log().Warn("节点重复上线（数据库检查）", attrs...)
		_ = as.writeError(sig, http.StatusConflict, nil)
		return nil, nil, false
	}

	brokerID := broker.ID
	peer := linkhub.NewPeer(brokerID, mux)
	if !as.opt.huber.Put(peer) {
		as.log().Warn("节点重复上线（连接池检查）", attrs...)
		_ = as.writeError(sig, http.StatusConflict, nil)
		return nil, nil, false
	}

	// TODO 响应成功报文
	initalCfg := authmesg.BrokerInitialConfig{URI: as.cfg.Database.URI}
	if err = as.writeSucceed(sig, initalCfg); err != nil {
		as.opt.huber.DelByID(brokerID)
		attrs = append(attrs, slog.Any("error", err))
		as.log().Warn("响应报文写入失败", attrs...)
		return nil, nil, false
	}

	// 修改数据库在线状态
	tunStat := &model.TunnelStat{
		ConnectedAt: now,
		KeepaliveAt: now,
		Protocol:    protocol,
		Subprotocol: subprotocol,
		LocalAddr:   mux.Addr().String(),
		RemoteAddr:  mux.RemoteAddr().String(),
	}
	exeStat := &model.ExecuteStat{
		Goos:       req.Goos,
		Goarch:     req.Goarch,
		PID:        req.PID,
		Args:       req.Args,
		Hostname:   req.Hostname,
		Workdir:    req.Workdir,
		Executable: req.Executable,
		Username:   req.Username,
		UID:        req.UID,
	}
	update := bson.M{"$set": bson.M{
		"status": true, "tunnel_stat": tunStat, "execute_stat": exeStat,
	}}
	filter := bson.M{"_id": brokerID, "status": false}
	brokerRepo := as.repo.Broker()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ret, err1 := brokerRepo.UpdateOne(ctx, filter, update)
	if err1 == nil && ret.ModifiedCount != 0 {
		as.log().Info("节点上线成功", attrs...)
		return req, peer, true
	}

	as.opt.huber.DelByID(brokerID)
	_ = as.writeError(sig, http.StatusInternalServerError, err1)

	if err1 != nil {
		attrs = append(attrs, slog.Any("error", err1))
	}
	as.log().Error("节点上线失败", attrs...)

	return nil, nil, false
}

// checkout 获得 agent 节点的信息，如果不存在自动创建。
func (as *brokerServer) checkout(req *authmesg.BrokerToServerRequest, timeout time.Duration) (*model.Broker, error) {
	brokerRepo := as.repo.Broker()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return brokerRepo.FindOne(ctx, bson.M{"secret": req.Secret})
}

func (as *brokerServer) disconnected(peer linkhub.Peer, timeout time.Duration) {
	now := time.Now()
	id := peer.ID()
	rx, tx := peer.Muxer().Transferred()
	update := bson.M{"$set": bson.M{
		"status": false, "tunnel_stat.disconnected_at": now,
		"tunnel_stat.receive_bytes": tx, "tunnel_stat.transmit_bytes": rx,
		// 注意：此时要以 broker 视角统计流量，所以 rx tx 要互换一下。
	}}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	brokerRepo := as.repo.Broker()
	_, _ = brokerRepo.UpdateByID(ctx, id, update)
	as.opt.huber.DelByID(id)
}

func (as *brokerServer) log() *slog.Logger {
	if l := as.opt.logger; l != nil {
		return l
	}

	return slog.Default()
}

func (as *brokerServer) getServer(p linkhub.Peer) *http.Server {
	srv := as.opt.server
	if srv == nil {
		srv = &http.Server{Handler: http.NotFoundHandler()}
	}
	baseCtxFunc := srv.BaseContext
	srv.BaseContext = func(ln net.Listener) context.Context {
		ctx := context.Background()
		if baseCtxFunc != nil {
			ctx = baseCtxFunc(ln)
		}

		return linkhub.WithValue(ctx, p)
	}

	return srv
}

func (as *brokerServer) validate(req *authmesg.BrokerToServerRequest) error {
	if v := as.opt.valid; v != nil {
		return v.Validate(req)
	}

	return nil
}

// allow 做简单的并发限制，在 broker 重启时，会导致大批量的 agent 重连，
// agent 上线的各种判断与状态修改是耗资源操作，可以通过限流器抑制。
func (as *brokerServer) allow() bool {
	if limit := as.opt.limit; limit != nil {
		return limit.Allow()
	}

	return true
}

func (as *brokerServer) writeError(w io.Writer, code int, err error) error {
	resp := &authmesg.ServerToBrokerResponse{Code: code}
	if err != nil {
		resp.Message = err.Error()
	}

	return tunutil.WriteHead(w, resp)
}

func (as *brokerServer) writeSucceed(w io.Writer, cfg authmesg.BrokerInitialConfig) error {
	resp := &authmesg.ServerToBrokerResponse{Code: http.StatusAccepted, Config: cfg}
	return tunutil.WriteHead(w, resp)
}
