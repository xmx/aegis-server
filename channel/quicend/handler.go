package quicend

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xmx/aegis-server/channel/transport"
	"github.com/xmx/aegis-server/datalayer/repository"
)

func New(repo repository.All, hub transport.Huber, next http.Handler, log *slog.Logger) transport.Handler {
	return &fakeHandler{
		repo: repo,
		hub:  hub,
		next: next,
		log:  log,
	}
}

type fakeHandler struct {
	repo repository.All
	hub  transport.Huber
	next http.Handler
	log  *slog.Logger
}

func (fh *fakeHandler) Handle(mux transport.Muxer) error {
	//goland:noinspection GoUnhandledErrorResult
	defer mux.Close()
	if err := fh.handshake(mux, 10*time.Second); err != nil {
		return err
	}

	return nil
}

func (fh *fakeHandler) handshake(mux transport.Muxer, timeout time.Duration) error {
	timer := time.AfterFunc(timeout, func() {
		_ = mux.Close()
		fh.log.Info("等待客户端超时")
	})

	stm, err := mux.Accept()
	if err != nil {
		return err
	}
	timer.Stop()
	//goland:noinspection GoUnhandledErrorResult
	defer stm.Close()

	now := time.Now()
	deadline := now.Add(timeout)
	_ = stm.SetDeadline(deadline)

	req := new(transport.AuthRequest)
	if _, err = req.ReadFrom(stm); err != nil {
		return err
	}

	attr := slog.Any("request", req)
	id := req.ID
	peer := &agentPeer{id: id, mux: mux}
	res := new(transport.AuthResponse)
	if absent := fh.hub.PutIfAbsent(peer); !absent {
		res.Message = "节点重复上线"
		_, _ = res.WriteTo(stm)
		fh.log.Warn("节点重复上线", attr)
		return errors.New("节点重复上线")
	}
	defer func() {
		fh.hub.Del(id)
		fh.log.Warn("节点下线了", attr)
	}()

	res.Succeed = true
	if _, err = res.WriteTo(stm); err != nil {
		fh.log.Error("报文写入失败", attr, slog.Any("error", err))
		return err
	}

	fh.log.Info("节点上线", attr)
	_ = stm.Close()

	srv := &http.Server{
		Handler: fh.next,
		BaseContext: func(ln net.Listener) context.Context {
			return transport.WithValue(context.Background(), peer)
		},
	}
	err = srv.Serve(mux)

	return err
}
