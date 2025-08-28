package gateway

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-server/channel/transport"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewServer(repo repository.All, hub transport.Huber, next http.Handler, log *slog.Logger) transport.Handler {
	return &brokerGateway{
		repo: repo,
		hub:  hub,
		next: next,
		log:  log,
	}
}

type brokerGateway struct {
	repo repository.All
	hub  transport.Huber
	next http.Handler
	log  *slog.Logger
}

//goland:noinspection GoUnhandledErrorResult
func (bg *brokerGateway) Handle(mux transport.Muxer) error {
	defer mux.Close()

	// 开始握手校验
	peer, err := bg.precheck(mux, 10*time.Second)
	if err != nil {
		return err
	}

	attrs := []any{slog.String("id", peer.ID())}
	defer func() {
		if exx := bg.disconnect(peer); exx != nil {
			attrs = append(attrs, slog.Any("error", exx))
			bg.log.Error("修改节点下线状态错误", attrs...)
		}
		bg.log.Info("节点下线", attrs...)
		// TODO 节点下线 Hook
	}()

	bg.log.Info("节点上线", attrs...)
	// TODO 上线成功 Hook

	srv := &http.Server{
		Handler: bg.next,
		BaseContext: func(net.Listener) context.Context {
			return transport.WithValue(context.Background(), peer)
		},
	}
	if err = srv.Serve(mux); err != nil {
		attrs = append(attrs, slog.Any("error", err))
		bg.log.Warn("internal http serve 错误", attrs...)
	}

	return nil
}

func (bg *brokerGateway) precheck(mux transport.Muxer, timeout time.Duration) (*brokerPeer, error) {
	timer := time.AfterFunc(timeout, func() {
		_ = mux.Close()
	})

	sig, err := mux.Accept()
	if err != nil {
		bg.log.Warn("等待握手连接错误", "error", err)
		return nil, err
	}
	timer.Stop()
	//goland:noinspection GoUnhandledErrorResult
	defer sig.Close()

	now := time.Now()
	deadline := now.Add(timeout)
	_ = sig.SetDeadline(deadline)

	var req transport.AuthRequest
	if _, err = req.ReadFrom(sig); err != nil {
		bg.log.Warn("读取握手报文错误", "error", err)
		return nil, err
	}

	id, secret := req.ID, req.Secret
	oid, _ := bson.ObjectIDFromHex(id)
	brk := bg.lookupByID(oid)
	if brk == nil {
		res := &transport.AuthResponse{Message: "节点不存在"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	if secret == "" || secret != brk.Secret {
		res := &transport.AuthResponse{Message: "密钥错误"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	if brk.Status {
		res := &transport.AuthResponse{Message: "节点重复上线"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}

	goos, goarch := req.Goos, req.Goarch
	peer := &brokerPeer{
		id:     oid,
		mux:    mux,
		goos:   goos,
		goarch: goarch,
	}
	if !bg.hub.PutIfAbsent(peer) {
		res := &transport.AuthResponse{Message: "节点重复上线"}
		_, _ = res.WriteTo(sig)
		return nil, res
	}
	attrs := []any{slog.Any("request", req)}
	succ := &transport.AuthResponse{Succeed: true}
	if _, err = succ.WriteTo(sig); err != nil {
		bg.hub.Del(id)
		attrs = append(attrs, slog.Any("error", err))
		bg.log.Warn("成功报文写入错误", attrs...)
		return nil, err
	}

	update := bson.M{"$set": bson.M{
		"status":       true,
		"goos":         goos,
		"goarch":       goarch,
		"network":      mux.Network(),
		"remote_addr":  mux.RemoteAddr().String(),
		"alive_at":     now,
		"connected_at": now,
	}}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	brkRepo := bg.repo.Broker()
	if _, err = brkRepo.UpdateByID(ctx, oid, update); err != nil {
		bg.hub.Del(id)
		res := &transport.AuthResponse{Message: "内部错误，节点上线失败"}
		_, _ = res.WriteTo(sig)
		attrs = append(attrs, slog.Any("error", err))
		bg.log.Error("修改数据库节点上线状态错误", attrs...)
		return nil, err
	}

	return peer, nil
}

func (bg *brokerGateway) lookupByID(id bson.ObjectID) *model.Broker {
	if id.IsZero() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brkRepo := bg.repo.Broker()
	brk, err := brkRepo.FindByID(ctx, id)
	if err == nil || errors.Is(err, mongo.ErrNoDocuments) {
		return brk
	}
	bg.log.Error("查询节点错误", slog.Any("id", id), slog.Any("error", err))

	return nil
}

func (bg *brokerGateway) disconnect(peer *brokerPeer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brkRepo := bg.repo.Broker()
	_, err := brkRepo.UpdateByID(ctx, peer.id, bson.M{"$set": bson.M{"status": false}})
	bg.hub.Del(peer.ID())

	return err
}
