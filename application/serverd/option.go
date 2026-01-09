package serverd

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/xmx/aegis-control/linkhub"
)

// AuthConfigLoader broker 上线认证通过后需要加载基础的启动配置。
type AuthConfigLoader interface {

	// LoadAuthConfig 加载配置参数。
	LoadAuthConfig(ctx context.Context) (*AuthConfig, error)
}

type Options struct {
	ConnectListener ConnectListener
	ConfigLoader    AuthConfigLoader
	Handler         http.Handler
	Huber           linkhub.Huber
	Validator       func(any) error // 认证报文参数校验器
	Logger          *slog.Logger
	Timeout         time.Duration
	Context         context.Context
}

// ConnectListener 节点上下线通知接口。
type ConnectListener interface {
	// OnConnection 节点上线事件通知。
	//
	// 该方法为同步阻塞方法。
	OnConnection(now time.Time, peer linkhub.Peer)

	// OnDisconnection 节点下线事件通知。
	//
	// 该方法为同步阻塞方法。
	OnDisconnection(now time.Time, info linkhub.Info)
}
