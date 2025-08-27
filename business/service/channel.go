package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/xmx/aegis-server/channel/linkhub"
	"github.com/xmx/aegis-server/channel/transport"
	"github.com/xmx/aegis-server/contract/request"
	"github.com/xmx/aegis-server/contract/response"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/ship"
)

func NewChannel(repo repository.All, hub linkhub.Huber, next transport.Server, log *slog.Logger) *Channel {
	return &Channel{
		repo: repo,
		hub:  hub,
		next: next,
		log:  log,
	}
}

type Channel struct {
	repo repository.All
	hub  linkhub.Huber
	next transport.Server
	log  *slog.Logger
}

func (chn *Channel) Open(w http.ResponseWriter, r *http.Request, req *request.ChannelOpen) error {
	jack, ok := w.(http.Hijacker)
	if !ok {
		return ship.ErrBadRequest
	}

	conn, _, err := jack.Hijack()
	if err != nil {
		return err
	}

	mux, err := transport.NewTCP(conn)
	if err != nil {
		_ = conn.Close()
		return err
	}

	peer := chn.newPeer(mux, req)
	if absent := chn.hub.PutIfAbsent(peer); !absent {
		_ = mux.Close()
		chn.log.Info("节点已经在线", slog.Any("broker", req))
		return nil
	}

	// 写入成功报文
	if err = chn.writeAccepted(conn, r); err != nil {
		return err
	}

	defer func() {
		_ = mux.Close()
		chn.hub.Del(req.ID)
		chn.log.Warn("节点下线了", slog.Any("broker", req), slog.Any("error", err))
	}()

	chn.log.Info("节点上线了", slog.Any("broker", req))
	err = chn.next.Serve(mux)

	return nil
}

func (chn *Channel) newPeer(mux transport.Muxer, req *request.ChannelOpen) *httpPeer {
	return &httpPeer{
		id:  req.ID,
		mux: mux,
	}
}

func (chn *Channel) writeAccepted(c net.Conn, r *http.Request) error {
	header := http.Header{
		ship.HeaderContentType: []string{ship.MIMEApplicationJSONCharsetUTF8},
	}
	res := &http.Response{
		StatusCode: http.StatusAccepted, // 默认规定 http.StatusAccepted 为成功状态码
		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,
		Header:     header,
		Request:    r,
	}

	data := &response.ChannelOpen{
		Succeed: true,
	}
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return err
	}

	res.Body = io.NopCloser(body)
	res.ContentLength = int64(body.Len())
	if err := res.Write(c); err != nil {
		_ = c.Close()
		return err
	}

	return nil
}

type httpPeer struct {
	id  string
	mux transport.Muxer
}

func (h *httpPeer) ID() string {
	return h.id
}

func (h *httpPeer) Muxer() transport.Muxer {
	return h.mux
}
