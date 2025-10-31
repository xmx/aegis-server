package restapi

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/profile"
	"github.com/xmx/aegis-control/mongodb"
	"github.com/xmx/aegis-server/applet/initialize/request"
	"github.com/xmx/aegis-server/config"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func NewInstall(results chan<- *config.Config) *Install {
	return &Install{
		results: results,
	}
}

type Install struct {
	results chan<- *config.Config
	doing   atomic.Bool
	done    bool
}

func (inst *Install) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/install").POST(inst.setup)
	return nil
}

func (inst *Install) setup(c *ship.Context) error {
	req := new(request.InstallSetup)
	if err := c.Bind(req); err != nil {
		return err
	}

	if !inst.doing.CompareAndSwap(false, true) {
		return ship.ErrBadRequest.Newf("正在初始化中")
	}
	defer inst.doing.Store(false)
	if inst.done {
		return ship.ErrBadRequest.Newf("已初始化完毕")
	}

	c.Infof("准备初始化")
	parent := c.Request().Context()
	addr := req.Server.Addr

	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		c.Errorf("监听地址输入不合法", "addr", addr, "error", err)
		return err
	}
	pkt, err := net.ListenPacket("udp", addr)
	if err != nil {
		c.Errorf("UDP 监听不可用", "addr", addr, "error", err)
		return err
	}
	_ = pkt.Close()
	var same bool
	val := parent.Value(http.LocalAddrContextKey)
	if adr, _ := val.(*net.TCPAddr); adr != nil {
		same = taddr.Port == adr.Port
	}
	if !same {
		ln, err := net.ListenTCP("tcp", taddr)
		if err != nil {
			c.Errorf("TCP 监听不可用", "addr", addr, "error", err)
			return err
		}
		_ = ln.Close()
	}

	uri := req.Database.URI
	pu, err := connstring.Parse(uri)
	if err != nil {
		return err
	}
	dbname := pu.Database
	if dbname == "" {
		return ship.ErrBadRequest.Newf("数据库连接地址缺少库名")
	}

	db, err := mongodb.Open(uri)
	if err != nil {
		return err
	}
	cli := db.Client()
	defer cli.Disconnect(parent)

	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	if err = cli.Ping(ctx, nil); err != nil {
		return ship.ErrBadRequest.Newf("连接数据库错误：%s", err)
	}

	cfg := req.Config()
	if err = profile.WriteFile(config.Filename, cfg); err != nil {
		return err
	}
	c.Infof("配置初始化并保存成功")
	inst.done = true
	inst.results <- cfg

	return nil
}
