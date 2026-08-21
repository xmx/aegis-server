package linkhub

import (
	"log/slog"
	"net/http"

	"github.com/xmx/muxconn"
)

type MUXHandler interface {
	HandleMUX(mux muxconn.Muxer) error
}

type NodeMUX struct {
	biz http.Handler // 业务
	log *slog.Logger
}

func NewNodeMUX(log *slog.Logger) *NodeMUX {
	return &NodeMUX{
		log: log,
	}
}

func (h *NodeMUX) HandleMUX(mux muxconn.Muxer) error {
	return nil
}

func (h *NodeMUX) serveHTTP(mux muxconn.Muxer) error {
	srv := &http.Server{Handler: h.biz}
	return srv.Serve(mux) // 阻塞运行直至连接断开
}

func (h *NodeMUX) onconnect() {
	// 新增/更新 agent 表
}

func (h *NodeMUX) ondisconnect() {
	// 更新 agent 表与连接历史表
}
