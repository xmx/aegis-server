package serverd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/xmx/aegis-common/muxlink/muxconn"
	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-common/muxlink/muxtool"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewServer(repo repository.All, opts Options) muxproto.MUXAccepter {
	return &centralServer{repo: repo, opts: opts}
}

type centralServer struct {
	repo repository.All
	opts Options
}

// AcceptMUX 接受 broker 节点建立的连接。
//
//goland:noinspection GoUnhandledErrorResult
func (ctl *centralServer) AcceptMUX(mux muxconn.Muxer) {
	defer mux.Close()

	connectAt := time.Now()
	peer, err := ctl.authentication(mux)
	if err != nil {
		raddr := mux.RemoteAddr()
		ctl.log().Warn("节点上线失败", "remote_addr", raddr, "error", err)
		return
	}

	info := peer.Info()
	ctl.log().Info("节点上线成功", "info", info)
	if sh := ctl.opts.ServerHooker; sh != nil {
		sh.OnConnected(info, connectAt)
	}

	err = ctl.serveHTTP(peer)
	ctl.log().Warn("节点下线了", "info", info, "error", err)

	ctl.disconnection(peer, connectAt)
}

//goland:noinspection GoUnhandledErrorResult
func (ctl *centralServer) authentication(mux muxconn.Muxer) (linkhub.Peer, error) {
	timeout := ctl.timeout()

	fc := muxtool.NewFlagCloser(mux)
	timer := time.AfterFunc(timeout, fc.Close)
	conn, err := mux.Accept()
	timer.Stop()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if fc.Closed() {
		return nil, net.ErrClosed
	}

	req := new(AuthRequest)
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	if err = muxtool.ReadAuth(conn, req); err != nil {
		ctl.log().Warn("读取认证报文错误", "error", err)
		return nil, err
	}

	if err = ctl.validAuthRequest(req); err != nil {
		req.Secret = "******" // 防止打印日志时泄露
		ctl.log().Warn("认证报文参数校验错误", "request", req, "error", err)
		ctl.responseError(conn, err, 0)
		return nil, err
	}

	secret := req.Secret
	req.Secret = "******" // 防止打印日志时泄露
	attrs := []any{"request", req}
	ctl.log().Debug("收到认证消息", attrs...)
	brok, err1 := ctl.findBrokerBySecret(secret)
	if err1 != nil {
		var code int
		if errors.Is(err, mongo.ErrNilDocument) {
			code = http.StatusNotFound
			err1 = errors.New("此节点尚未注册")
		} else {
			err1 = fmt.Errorf("查询节点出错：%w", err1)
		}

		attrs = append(attrs, "error", err1)
		ctl.log().Warn("认证报文参数校验错误", attrs...)
		ctl.responseError(conn, err1, code)

		return nil, err
	}

	// 在线状态检查
	if brok.Status {
		err = errors.New("此节点已经在线了（数据库）")
		ctl.log().Warn("节点重复上线（数据库）", attrs...)
		ctl.responseError(conn, err, http.StatusConflict)

		return nil, err
	}

	brokID := brok.ID
	info := linkhub.Info{
		Name: brok.Name, Inet: req.Inet, Goos: req.Goos, Goarch: req.Goarch,
		Hostname: req.Hostname, Semver: req.Semver,
	}
	peer := ctl.putHuber(brokID, mux, info)
	if peer == nil {
		err = errors.New("此节点已经在线了（连接池）")
		ctl.log().Warn("节点重复上线（连接池）", attrs...)
		ctl.responseError(conn, err, http.StatusConflict)

		return nil, err
	}

	cfg, err := ctl.loadAuthConfig()
	if err != nil {
		ctl.removeHuber(brokID) // 获取配置错误，从连接池中删除并返回错误。

		attrs = append(attrs, "error", err)
		ctl.log().Error("加载初始配置文件出错", attrs...)
		ctl.responseError(conn, err, 0)
		return nil, err
	}

	if err = ctl.responseConfig(conn, cfg); err != nil {
		ctl.removeHuber(brokID) // 返回消息出错，从连接池中删除并返回错误。

		attrs = append(attrs, "error", err)
		ctl.log().Error("认证响应消息返回出错", attrs...)

		return nil, err
	}

	if ret, err2 := ctl.updateBrokerOnline(mux, req, brok); err2 != nil || ret.ModifiedCount == 0 {
		ctl.removeHuber(brokID) // 修改数据库状态失败，从连接池中删除并返回错误。

		if err2 == nil {
			err2 = errors.New("没有找到该节点（修改在线状态）")
		}
		attrs = append(attrs, "error", err)
		ctl.log().Error("节点重复上线（连接池）", attrs...)
		ctl.responseError(conn, err2, http.StatusConflict)

		return nil, err
	}

	return peer, nil
}

func (ctl *centralServer) log() *slog.Logger {
	if l := ctl.opts.Logger; l != nil {
		return l
	}
	return slog.Default()
}

func (ctl *centralServer) serveHTTP(peer linkhub.Peer) error {
	h := ctl.opts.Handler
	if h == nil {
		h = http.NotFoundHandler()
	}

	srv := &http.Server{
		Handler: h,
		BaseContext: func(net.Listener) context.Context {
			return linkhub.WithValue(context.Background(), peer)
		},
	}
	mux := peer.Muxer()

	return srv.Serve(mux)
}

func (ctl *centralServer) disconnection(peer linkhub.Peer, connectAt time.Time) {
	disconnectAt := time.Now()
	id := peer.ID()
	info := peer.Info()
	mux := peer.Muxer()
	tx, rx := mux.Traffic() // 互换

	attrs := []any{"info", info}
	filter := bson.D{{"_id", id}, {"status", true}}
	update := bson.M{"$set": bson.M{
		"status": false, "tunnel_stat.disconnected_at": disconnectAt,
		"tunnel_stat.receive_bytes": rx, "tunnel_stat.transmit_bytes": tx,
		// 注意：此时是站在 broker 视角统计的流量，所以 rx tx 要互换一下。
	}}

	ctx, cancel := ctl.perContext()
	defer cancel()

	repo := ctl.repo.Broker()
	if ret, err := repo.UpdateOne(ctx, filter, update); err != nil {
		attrs = append(attrs, "error", err)
		ctl.log().Error("修改数据库节点下线状态错误", attrs...)
	} else if ret.ModifiedCount == 0 {
		ctl.log().Error("修改数据库节点下线状态无修改", attrs...)
	}

	ctl.removeHuber(id)

	libName, libModule := mux.Library()
	raddr, laddr := mux.Addr(), mux.RemoteAddr() // 互换
	second := int64(disconnectAt.Sub(connectAt).Seconds())
	history := &model.BrokerConnectHistory{
		Broker: id,
		Name:   info.Name,
		Semver: info.Semver,
		Inet:   info.Inet,
		Goos:   info.Goos,
		Goarch: info.Goarch,
		TunnelStat: model.TunnelStatHistory{
			ConnectedAt:    connectAt,
			DisconnectedAt: disconnectAt,
			Second:         second,
			Library:        model.TunnelLibrary{Name: libName, Module: libModule},
			LocalAddr:      laddr.String(),
			RemoteAddr:     raddr.String(),
			ReceiveBytes:   rx,
			TransmitBytes:  tx,
		},
	}
	hisRepo := ctl.repo.BrokerConnectHistory()
	if _, err := hisRepo.InsertOne(ctx, history); err != nil {
		attrs = append(attrs, "save_history_error", err)
		ctl.log().Error("保存连接历史记录错误", attrs...)
	}

	ctl.log().Info("节点下线处理完毕", attrs...)

	if sh := ctl.opts.ServerHooker; sh != nil {
		sh.OnDisconnected(info, connectAt, disconnectAt)
	}
}

func (ctl *centralServer) timeout() time.Duration {
	if du := ctl.opts.Timeout; du > 0 {
		return du
	}

	return time.Minute
}

func (ctl *centralServer) perContext() (context.Context, context.CancelFunc) {
	d := ctl.timeout()
	return context.WithTimeout(ctl.opts.Context, d)
}

func (ctl *centralServer) loadAuthConfig() (*AuthConfig, error) {
	cl := ctl.opts.ConfigLoader
	ctx, cancel := ctl.perContext()
	defer cancel()

	return cl.LoadAuthConfig(ctx)
}

func (ctl *centralServer) validAuthRequest(req *AuthRequest) error {
	if f := ctl.opts.Validator; f != nil {
		return f(req)
	}

	var errs []error
	if req.Secret == "" {
		errs = append(errs, errors.New("连接密钥必须填写 (secret)"))
	}
	if _, err := netip.ParseAddr(req.Inet); err != nil {
		errs = append(errs, errors.New("出口网卡地址必须填写 (inet)"))
	}
	if req.Goos == "" {
		errs = append(errs, errors.New("操作系统类型必须填写 (goos)"))
	}
	if req.Goarch == "" {
		errs = append(errs, errors.New("架构必须填写 (goarch)"))
	}
	if req.Semver == "" {
		errs = append(errs, errors.New("版本号必须填写 (semver)"))
	}

	return errors.Join(errs...)
}

func (ctl *centralServer) responseError(conn net.Conn, err error, code int) error {
	if code < http.StatusBadRequest {
		code = http.StatusBadRequest
	}
	dat := &authResponse{Code: code, Message: err.Error()}

	d := ctl.timeout()
	_ = conn.SetWriteDeadline(time.Now().Add(d))

	return muxtool.WriteAuth(conn, dat)
}

func (ctl *centralServer) responseConfig(conn net.Conn, cfg *AuthConfig) error {
	dat := &authResponse{Code: http.StatusAccepted, Config: cfg}

	d := ctl.timeout()
	_ = conn.SetWriteDeadline(time.Now().Add(d))

	return muxtool.WriteAuth(conn, dat)
}

func (ctl *centralServer) findBrokerBySecret(secret string) (*model.Broker, error) {
	d := ctl.timeout()
	ctx, cancel := context.WithTimeout(ctl.opts.Context, d)
	defer cancel()

	repo := ctl.repo.Broker()
	opt := options.FindOne().SetProjection(bson.M{"_id": 1, "status": 1, "name": 1})

	return repo.FindOne(ctx, bson.D{{"secret", secret}}, opt)
}

func (ctl *centralServer) updateBrokerOnline(mux muxconn.Muxer, req *AuthRequest, before *model.Broker) (*mongo.UpdateResult, error) {
	now := time.Now()
	libName, libModule := mux.Library()
	raddr, laddr := mux.Addr(), mux.RemoteAddr() // 位置互换
	tunStat := &model.TunnelStat{
		ConnectedAt: now,
		KeepaliveAt: now,
		Library:     model.TunnelLibrary{Name: libName, Module: libModule},
		LocalAddr:   raddr.String(),
		RemoteAddr:  laddr.String(),
	}
	exeStat := &model.ExecuteStat{
		Inet:       req.Inet,
		Goos:       req.Goos,
		Goarch:     req.Goarch,
		Semver:     req.Semver,
		PID:        req.PID,
		Args:       req.Args,
		Hostname:   req.Hostname,
		Workdir:    req.Workdir,
		Executable: req.Executable,
	}

	update := bson.M{"$set": bson.M{
		"status": true, "tunnel_stat": tunStat, "execute_stat": exeStat,
	}}
	filter := bson.D{{"_id", before.ID}, {"status", false}}

	d := ctl.timeout()
	ctx, cancel := context.WithTimeout(ctl.opts.Context, d)
	defer cancel()

	repo := ctl.repo.Broker()

	return repo.UpdateOne(ctx, filter, update)
}

func (ctl *centralServer) putHuber(id bson.ObjectID, mux muxconn.Muxer, inf linkhub.Info) linkhub.Peer {
	return ctl.opts.Huber.Put(id, mux, inf)
}

func (ctl *centralServer) removeHuber(id bson.ObjectID) {
	ctl.opts.Huber.DelID(id)
}
