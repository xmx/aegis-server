package initapi

import (
	"context"
	"net"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/agies-server/argument/request"
	"github.com/xmx/agies-server/handler/shipx"
	"github.com/xmx/agies-server/library/sqldb"
)

func Testing() shipx.Register {
	return &testingAPI{}
}

type testingAPI struct{}

func (api *testingAPI) Register(route shipx.Router) error {
	route.Anon().Route("/testing/listen").POST(api.Listen)
	route.Anon().Route("/testing/tidb").POST(api.TiDB)
	return nil
}

// Listen 测试监听本地某个地址。
func (api *testingAPI) Listen(c *ship.Context) error {
	req := new(request.TestingListen)
	if err := c.Bind(req); err != nil {
		return err
	}

	parent := c.Request().Context()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	lc := new(net.ListenConfig)
	lis, err := lc.Listen(ctx, "tcp", req.Addr)
	if err != nil {
		return err
	}
	_ = lis.Close()

	return nil
}

// TiDB 测试连接数据库。
func (api *testingAPI) TiDB(c *ship.Context) error {
	req := new(request.TestingTiDB)
	if err := c.Bind(req); err != nil {
		return err
	}

	db, err := sqldb.TiDB(req.DSN, 5*time.Second)
	if err != nil {
		return err
	}
	_ = db.Close()

	return nil
}
