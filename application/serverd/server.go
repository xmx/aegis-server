package serverd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-common/options"
	"github.com/xmx/aegis-common/tunnel/tunconst"
	"github.com/xmx/aegis-common/tunnel/tunopen"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"github.com/xmx/aegis-server/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func New(repo repository.All, cfg *config.Config, opts ...options.Lister[option]) tunconst.Handler {
	opts = append(opts, fallbackOptions())
	opt := options.Eval[option](opts...)

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

func (bs *brokerServer) Handle(mux tunopen.Muxer) {
	//goland:noinspection GoUnhandledErrorResult
	defer mux.Close()

	if !bs.opt.allow() {
		bs.log().Warn("限流器抑制 broker 上线")
		return
	}

	timeout := bs.opt.timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	peer, succeed := bs.authentication(mux, timeout)
	if !succeed {
		return
	}
	defer bs.disconnected(peer, timeout)

	srv := bs.getServer(peer)
	_ = srv.Serve(mux)
}

// authentication 节点认证。
// 客户端主动建立一条虚拟子流连接用于交换认证信息，认证后改子流关闭，后续子流即为业务流。
func (bs *brokerServer) authentication(mux tunopen.Muxer, timeout time.Duration) (linkhub.Peer, bool) {
	protocol, subprotocol := mux.Protocol()
	laddr, raddr := mux.Addr(), mux.RemoteAddr()
	attrs := []any{
		slog.String("protocol", protocol), slog.String("subprotocol", subprotocol),
		slog.Any("local_addr", laddr), slog.Any("remote_addr", raddr),
	}

	// 设置超时主动断开，防止恶意客户端一直不建立认证连接。
	timer := time.AfterFunc(timeout, func() { _ = mux.Close() })
	defer timer.Stop()

	sig, err := mux.Accept()
	timer.Stop()
	if err != nil {
		bs.log().Error("等待客户端建立认证连接出错", "error", err)
		return nil, false
	}
	//goland:noinspection GoUnhandledErrorResult
	defer sig.Close()

	// 读取数据
	now := time.Now()
	_ = sig.SetDeadline(now.Add(timeout))
	req := new(authRequest)
	if err = tunopen.ReadAuth(sig, req); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		bs.log().Error("读取请求信息错误", attrs...)
		return nil, false
	}

	attrs = append(attrs, slog.Any("auth_request", req))
	if err = bs.opt.valid(req); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		bs.log().Error("读取请求信息校验错误", attrs...)
		_ = bs.writeError(sig, http.StatusBadRequest, err)
		return nil, false
	}
	brok, err := bs.findBroker(req.Secret, timeout)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNilDocument) {
			bs.log().Warn("broker 节点不存在")
			code = http.StatusNotFound
		} else {
			attrs = append(attrs, slog.Any("error", err))
			bs.log().Error("查询 broker 节点错误", attrs...)
		}

		_ = bs.writeError(sig, code, err)
		return nil, false
	}
	if brok.Status { // 节点已经在线了
		bs.log().Warn("节点重复上线（数据库检查）", attrs...)
		_ = bs.writeError(sig, http.StatusConflict, nil)
		return nil, false
	}

	brokerID := brok.ID
	pinf := linkhub.Info{Inet: req.Inet, Goos: req.Goos, Goarch: req.Goarch, Hostname: req.Hostname}
	peer := linkhub.NewPeer(brokerID, mux, pinf)
	if !bs.opt.huber.Put(peer) {
		bs.log().Warn("节点重复上线（连接池检查）", attrs...)
		_ = bs.writeError(sig, http.StatusConflict, nil)
		return nil, false
	}

	authCfg := authConfig{URI: bs.cfg.Database.URI}
	if err = bs.writeSucceed(sig, authCfg); err != nil {
		bs.opt.huber.DelByID(brokerID)
		attrs = append(attrs, slog.Any("error", err))
		bs.log().Warn("响应报文写入失败", attrs...)
		return nil, false
	}

	// 修改数据库在线状态
	tunStat := &model.TunnelStat{
		ConnectedAt: now,
		KeepaliveAt: now,
		Protocol:    protocol,
		Subprotocol: subprotocol,
		LocalAddr:   raddr.String(), // 位置互换
		RemoteAddr:  laddr.String(),
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
	}
	update := bson.M{"$set": bson.M{
		"status": true, "tunnel_stat": tunStat, "execute_stat": exeStat,
	}}
	filter := bson.D{{"_id", brokerID}, {"status", false}}
	brokerRepo := bs.repo.Broker()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ret, err1 := brokerRepo.UpdateOne(ctx, filter, update)
	if err1 == nil && ret.ModifiedCount != 0 {
		bs.log().Info("节点上线成功", attrs...)
		return peer, true
	}

	bs.opt.huber.DelByID(brokerID)
	_ = bs.writeError(sig, http.StatusInternalServerError, err1)

	if err1 != nil {
		attrs = append(attrs, slog.Any("error", err1))
	}
	bs.log().Error("节点上线失败", attrs...)

	return nil, false
}

func (bs *brokerServer) findBroker(secret string, timeout time.Duration) (*model.Broker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	repo := bs.repo.Broker()

	return repo.FindOne(ctx, bson.M{"secret": secret})
}

func (bs *brokerServer) disconnected(peer linkhub.Peer, timeout time.Duration) {
	now := time.Now()
	id := peer.ID()
	rx, tx := peer.Muxer().Traffic()
	update := bson.M{"$set": bson.M{
		"status": false, "tunnel_stat.disconnected_at": now,
		"tunnel_stat.receive_bytes": tx, "tunnel_stat.transmit_bytes": rx,
		// 注意：此时是站在 broker 视角统计的流量，所以 rx tx 要互换一下。
	}}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	brokerRepo := bs.repo.Broker()
	_, _ = brokerRepo.UpdateByID(ctx, id, update)
	bs.opt.huber.DelByID(id)
}

func (bs *brokerServer) log() *slog.Logger {
	if l := bs.opt.logger; l != nil {
		return l
	}

	return slog.Default()
}

func (bs *brokerServer) getServer(p linkhub.Peer) *http.Server {
	srv := bs.opt.server
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

func (bs *brokerServer) writeError(w io.Writer, code int, err error) error {
	resp := &authResponse{Code: code}
	if err != nil {
		resp.Message = err.Error()
	}

	return tunopen.WriteAuth(w, resp)
}

func (bs *brokerServer) writeSucceed(w io.Writer, cfg authConfig) error {
	resp := &authResponse{Code: http.StatusAccepted, Config: cfg}
	return tunopen.WriteAuth(w, resp)
}
