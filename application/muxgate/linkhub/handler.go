package linkhub

import (
	"log/slog"
	"net/http"

	"github.com/xmx/muxconn"
)

type MUXHandler interface {
	HandleMUX(mux muxconn.Muxer) error
}

type AgentMUX struct {
	biz http.Handler // 业务
	log *slog.Logger
}

func NewAgentMUX(log *slog.Logger) *AgentMUX {
	return &AgentMUX{
		log: log,
	}
}

func (h *AgentMUX) HandleMUX(mux muxconn.Muxer) error {
	return h.serveHTTP(mux)
}

func (h *AgentMUX) serveHTTP(mux muxconn.Muxer) error {
	srv := &http.Server{Handler: h.biz}
	return srv.Serve(mux) // 阻塞运行直至连接断开
}

func (h *AgentMUX) onconnect() {
	// 新增/更新 agent 表
}

func (h *AgentMUX) ondisconnect() {
	// 更新 agent 表与连接历史表
}
