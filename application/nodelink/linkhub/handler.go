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

func (s *AgentMUX) HandleMUX(mux muxconn.Muxer) error {
	return s.serveHTTP(mux)
}

func (s *AgentMUX) serveHTTP(mux muxconn.Muxer) error {
	srv := &http.Server{Handler: s.biz}
	return srv.Serve(mux) // 阻塞运行直至连接断开
}

func (s *AgentMUX) onconnect() {
	// 新增/更新 agent 表
}

func (s *AgentMUX) ondisconnect() {
	// 更新 agent 表与连接历史表
}
