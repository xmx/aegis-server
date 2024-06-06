package initapi

import (
	"net"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/agies-server/argument/request"
)

func Testing() {

}

type testingAPI struct{}

// Listen 测试监听本地某个地址。
func (api *testingAPI) Listen(c *ship.Context) error {
	req := new(request.TestingListen)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	lc := new(net.ListenConfig)
	lis, err := lc.Listen(ctx, "tcp", req.Addr)
	if err != nil {
		return err
	}
	_ = lis.Close()

	return nil
}

// DB 测试连接数据库。
func (api *testingAPI) DB(c *ship.Context) error {
	return nil
}
