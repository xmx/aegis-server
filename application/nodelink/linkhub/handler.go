package linkhub

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-common/muxutil"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/muxconn"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MUXHandler interface {
	HandleMUX(mux muxconn.Muxer) error
}

type AgentMUX struct {
	db  *repository.BaseDB // 数据库
	biz http.Handler       // 业务
	log *slog.Logger       // 日志
}

func NewAgentMUX(db *repository.BaseDB, biz http.Handler, log *slog.Logger) *AgentMUX {
	return &AgentMUX{
		db:  db,
		biz: biz,
		log: log,
	}
}

func (s *AgentMUX) HandleMUX(mux muxconn.Muxer) error {
	// 通道建立并获取认证报文
	connAt := time.Now()
	req, dat, err := s.handshake(mux, connAt)
	if err != nil {
		return err
	}

	cli := newClientMUX(mux, req, dat)
	defer func() {
		_ = cli.Close()
		_ = s.updateAgentOfflineDB(cli, connAt)
	}()

	return s.serveHTTP(cli)
}

func (s *AgentMUX) handshake(mux muxconn.Muxer, connAt time.Time) (*muxutil.Request, *model.Agent, error) {
	// 10s 之内 agent 必须要建立认证连接，否则认为协议不匹配或恶意连接。
	sc := newSafeClose(mux)
	timer := time.AfterFunc(10*time.Second, sc.close)
	conn, err := mux.Accept()
	timer.Stop()

	if err != nil {
		return nil, nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()
	if sc.closed() {
		return nil, nil, context.DeadlineExceeded
	}

	// 读取 agent 发来的身份信息。
	req := new(muxutil.Request)
	if err = muxutil.ReadAuth(conn, req); err != nil {
		return nil, nil, err
	}

	dat, err := s.findOrCreate(req)
	if err != nil {
		msg := muxutil.NewErrResponse[any](http.StatusBadRequest, err)
		_ = muxutil.WriteAuth(conn, msg)

		return nil, nil, err
	}
	if !dat.Enabled {
		msg := muxutil.NewMsgResponse[any](http.StatusForbidden, "禁止登录")
		_ = muxutil.WriteAuth(conn, msg)

		return nil, nil, msg
	}
	if dat.Status {
		msg := muxutil.NewMsgResponse[any](http.StatusConflict, "节点已经在线")
		_ = muxutil.WriteAuth(conn, msg)

		return nil, nil, msg
	}

	if msg := s.updateAgentOnlineDB(dat.ID, mux, req, connAt); msg != nil {
		_ = muxutil.WriteAuth(conn, msg)
		return nil, nil, msg
	}

	// TODO put hub

	msg := muxutil.NewOKResponse[any](nil)
	if err = muxutil.WriteAuth(conn, msg); err != nil {
		return nil, nil, err
	}

	return req, dat, nil
}

func (s *AgentMUX) findOrCreate(req *muxutil.Request) (*model.Agent, error) {
	machineID := req.MachineID
	ctx, cancel := s.perContext()
	defer cancel()

	coll := s.db.Agent()
	if dat, err := coll.FindOne(ctx, bson.D{{Key: "machine_id", Value: machineID}}); err == nil || !errors.Is(err, mongo.ErrNoDocuments) {
		return dat, err
	}

	now := time.Now()
	dat := &model.Agent{
		MachineID: machineID,
		Enabled:   true, // 默认启用
		CreatedAt: now,
		UpdatedAt: now,
	}

	ret, err := coll.InsertOne(ctx, dat)
	if err != nil {
		return nil, err
	}
	dat.ID, _ = ret.InsertedID.(bson.ObjectID)

	return dat, nil
}

func (s *AgentMUX) serveHTTP(cli *clientMUX) error {
	next := s.biz
	if next == nil {
		next = http.NotFoundHandler()
	}

	srv := &http.Server{
		Handler: next,
		BaseContext: func(net.Listener) context.Context {
			return withContext(cli)
		},
	}

	return srv.Serve(cli) // 阻塞运行直至连接断开
}

func (s *AgentMUX) updateAgentOnlineDB(id bson.ObjectID, mux muxconn.Muxer, req *muxutil.Request, connAt time.Time) *muxutil.Response[any] {
	semver := model.ParseSemver(req.Version)
	proc := &model.AgentProcessInfo{
		Inet:       req.Inet,
		Semver:     semver,
		GOOS:       req.GOOS,
		GOARCH:     req.GOARCH,
		PID:        req.PID,
		Args:       req.Args,
		Hostname:   req.Hostname,
		Workdir:    req.Workdir,
		Executable: req.Executable,
		Environ:    req.Environ,
	}

	name, module := mux.Library()
	tun := &model.AgentTunnelInfo{
		ConnectedAt:   connAt,
		KeepaliveAt:   connAt,
		LibraryName:   name,
		LibraryModule: module,
		ServerAddr:    mux.Addr().String(),
		RemoteAddr:    mux.RemoteAddr().String(),
	}

	filter := bson.D{{Key: "_id", Value: id}, {Key: "status", Value: false}, {Key: "enabled", Value: true}}
	update := bson.M{"$set": bson.M{"status": true, "process_info": proc, "tunnel_info": tun, "updated_at": connAt}}

	ctx, cancel := s.perContext()
	defer cancel()

	coll := s.db.Agent()
	ret, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return muxutil.NewErrResponse[any](http.StatusBadRequest, err)
	}
	if ret.ModifiedCount == 0 {
		return muxutil.NewMsgResponse[any](http.StatusNotFound, "数据不存在")
	}

	return nil
}

func (s *AgentMUX) updateAgentOfflineDB(cli *clientMUX, connAt time.Time) error {
	disconnAt := time.Now()

	name, module := cli.Library()
	tx, rx := cli.Traffic() // 视角互换

	tun := model.AgentTunnelInfo{
		ConnectedAt:   connAt,
		KeepaliveAt:   disconnAt,
		LibraryName:   name,
		LibraryModule: module,
		ServerAddr:    cli.Addr().String(),
		RemoteAddr:    cli.RemoteAddr().String(),
		ReceiveBytes:  rx,
		TransmitBytes: tx,
	}

	inf := cli.Info()
	proc := model.AgentProcessInfo{
		Inet:       inf.Inet,
		Semver:     inf.Semver,
		GOOS:       inf.GOOS,
		GOARCH:     inf.GOARCH,
		PID:        inf.PID,
		Args:       inf.Args,
		Hostname:   inf.Hostname,
		Workdir:    inf.Workdir,
		Executable: inf.Executable,
		Environ:    inf.Environ,
	}

	id := inf.ID
	filter := bson.D{{Key: "_id", Value: id}, {Key: "status", Value: true}}
	update := bson.M{"$set": bson.M{"status": false, "tunnel_info": tun, "updated_at": disconnAt}}

	ctx, cancel := s.perContext()
	defer cancel()

	{
		coll := s.db.Agent()
		_, _ = coll.UpdateOne(ctx, filter, update)
	}

	his := &model.AgentConnRecord{
		AgentID:        id,
		MachineID:      inf.MachineID,
		ProcessInfo:    proc,
		TunnelInfo:     tun,
		ActiveSeconds:  int64(disconnAt.Sub(connAt).Seconds()),
		DisconnectedAt: disconnAt,
	}
	{
		coll := s.db.AgentConnRecord()
		_, _ = coll.InsertOne(ctx, his)
	}

	return nil
}

func (s *AgentMUX) perContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
