package broker

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/contract/linkhub"
	"github.com/xmx/aegis-control/contract/sbmesg"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewGate(repo repository.All, hub linkhub.Huber, next http.Handler, log *slog.Logger) transport.Handler {
	return &gateway{
		repo: repo,
		hub:  hub,
		next: next,
		log:  log,
	}
}

type gateway struct {
	repo repository.All
	hub  linkhub.Huber
	next http.Handler
	log  *slog.Logger
}

//goland:noinspection GoUnhandledErrorResult
func (gw *gateway) Handle(mux transport.Muxer) error {
	defer mux.Close()

	// 开始握手校验
	info, err := gw.precheck(mux, 10*time.Second)
	if err != nil {
		return err
	}

	proto := mux.Protocol()
	remoteAddr := mux.RemoteAddr()
	peer := info.peer
	id := peer.Host()
	attrs := []any{
		slog.String("id", id),
		slog.String("name", info.data.Name),
		slog.String("goos", info.req.Goos),
		slog.String("goarch", info.req.Goarch),
		slog.String("protocol", proto),
		slog.String("remote_addr", remoteAddr.String()),
	}
	defer func() {
		if exx := gw.disconnect(peer); exx != nil {
			attrs = append(attrs, slog.Any("error", exx))
			gw.log.Error("修改节点下线状态错误", attrs...)
		}
		gw.log.Info("节点下线", attrs...)
		// TODO 节点下线 Hook
	}()

	gw.log.Info("节点上线", attrs...)
	// TODO 上线成功 Hook

	srv := gw.newServer(peer)
	if err = srv.Serve(mux); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		gw.log.Warn("internal http serve 错误", attrs...)
	}

	return nil
}

func (gw *gateway) newServer(p linkhub.Peer) *http.Server {
	bg := context.Background()
	return &http.Server{
		Handler: gw.next,
		BaseContext: func(net.Listener) context.Context {
			return linkhub.WithValue(bg, p)
		},
	}
}

func (gw *gateway) precheck(mux transport.Muxer, timeout time.Duration) (*authInfo, error) {
	timer := time.AfterFunc(timeout, func() {
		_ = mux.Close()
	})

	sig, err := mux.Accept()
	if err != nil {
		gw.log.Warn("等待握手连接错误", "error", err)
		return nil, err
	}
	timer.Stop()
	//goland:noinspection GoUnhandledErrorResult
	defer sig.Close()

	now := time.Now()
	deadline := now.Add(timeout)
	_ = sig.SetDeadline(deadline)

	var req sbmesg.BrokerRequest
	if _, err = req.ReadFrom(sig); err != nil {
		gw.log.Warn("读取握手报文错误", "error", err)
		return nil, err
	}

	id, secret := req.ID, req.Secret
	oid, _ := bson.ObjectIDFromHex(id)
	brk := gw.lookupByID(oid)
	if brk == nil {
		res := &sbmesg.BrokerResponse{Message: "节点不存在"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	if secret == "" || secret != brk.Secret {
		res := &sbmesg.BrokerResponse{Message: "密钥错误"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	if brk.Status {
		res := &sbmesg.BrokerResponse{Message: "节点重复上线"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}

	peer := linkhub.NewPeer(oid, mux)
	if !gw.hub.Put(peer) {
		res := &sbmesg.BrokerResponse{Message: "节点重复上线"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	attrs := []any{slog.Any("request", req)}
	succ := &sbmesg.BrokerResponse{Succeed: true}
	if _, err = succ.WriteTo(sig); err != nil {
		gw.hub.Del(id)
		attrs = append(attrs, slog.Any("error", err))
		gw.log.Warn("成功报文写入错误", attrs...)
		return nil, err
	}

	raddr := mux.RemoteAddr().String()
	goos, goarch := req.Goos, req.Goarch
	update := bson.M{"$set": bson.M{
		"status":       true,
		"goos":         goos,
		"goarch":       goarch,
		"protocol":     mux.Protocol(),
		"remote_addr":  raddr,
		"alive_at":     now,
		"connected_at": now,
	}}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	brkRepo := gw.repo.Broker()
	if _, err = brkRepo.UpdateByID(ctx, oid, update); err != nil {
		gw.hub.Del(id)
		res := &sbmesg.BrokerResponse{Message: "内部错误，节点上线失败"}
		_, _ = res.WriteTo(sig)
		attrs = append(attrs, slog.Any("error", err))
		gw.log.Error("修改数据库节点上线状态错误", attrs...)
		return nil, err
	}
	info := &authInfo{
		peer: peer,
		data: brk,
		req:  req,
	}

	return info, nil
}

func (gw *gateway) lookupByID(id bson.ObjectID) *model.Broker {
	if id.IsZero() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brkRepo := gw.repo.Broker()
	brk, err := brkRepo.FindByID(ctx, id)
	if err == nil || errors.Is(err, mongo.ErrNoDocuments) {
		return brk
	}
	gw.log.Error("查询节点错误", slog.Any("id", id), slog.Any("error", err))

	return nil
}

func (gw *gateway) disconnect(p linkhub.Peer) error {
	now := time.Now()
	host, oid := p.Host(), p.ObjectID()
	update := bson.M{"$set": bson.M{"status": false, "disconnected_at": now}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brkRepo := gw.repo.Broker()
	_, err := brkRepo.UpdateByID(ctx, oid, update)
	gw.hub.Del(host)

	return err
}

type authInfo struct {
	peer linkhub.Peer
	data *model.Broker
	req  sbmesg.BrokerRequest
}
