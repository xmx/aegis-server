package initapi

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/library/sqldb"
)

// Testing 敏感接口。
func Testing() shipx.Register {
	return &testingAPI{}
}

type testingAPI struct{}

func (api *testingAPI) Register(route shipx.Router) error {
	route.Anon().Route("/testing/listen").POST(api.Listen)
	route.Anon().Route("/testing/tidb").POST(api.TiDB)
	route.Anon().Route("/testing/cert").POST(api.Cert)
	route.Anon().Route("/testing/config").POST(api.Config)
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
	if err == nil {
		_ = lis.Close()
		return nil
	}

	var listenPort, localPort int
	{
		ope, ok := err.(*net.OpError)
		if !ok {
			return err
		}
		addr, yes := ope.Addr.(*net.TCPAddr)
		if !yes {
			return err
		}
		listenPort = addr.Port
	}
	{
		addr, ok := parent.Value(http.LocalAddrContextKey).(*net.TCPAddr)
		if !ok {
			return err
		}
		localPort = addr.Port
	}
	if listenPort == localPort {
		err = nil
	}

	return err
}

// TiDB 测试连接数据库。
func (api *testingAPI) TiDB(c *ship.Context) error {
	req := new(request.TestingTiDB)
	if err := c.Bind(req); err != nil {
		return err
	}

	db, err := sqldb.TiDB(req.DSN)
	if err != nil {
		return err
	}

	parent := c.Request().Context()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer func() {
		cancel()
		_ = db.Close()
	}()

	var version string
	if err = db.QueryRowContext(ctx, "SELECT version()").
		Scan(&version); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, version)
}

func (api *testingAPI) Cert(c *ship.Context) error {
	req := new(request.TestingCert)
	if err := c.Bind(req); err != nil {
		return err
	}

	_, err := tls.X509KeyPair([]byte(req.Cert), []byte(req.Pkey))

	return err
}

func (api *testingAPI) Config(c *ship.Context) error {
	req := new(request.TestingConfig)
	if err := c.Bind(req); err != nil {
		return err
	}

	db, err := sqldb.TiDB(req.TiDB.DSN)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer db.Close()

	return nil
}
