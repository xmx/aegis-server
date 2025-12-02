package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/response"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewFirewall(svc *service.Firewall) *Firewall {
	return &Firewall{
		svc: svc,
	}
}

type Firewall struct {
	svc *service.Firewall
}

func (fw *Firewall) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/firewalls").GET(fw.list)
	r.Route("/firewall").
		POST(fw.create).
		PUT(fw.update).
		DELETE(fw.delete)
	r.Route("/firewall/precheck").POST(fw.precheck)
	r.Route("/firewall/reset").DELETE(fw.reset)

	return nil
}

// create 创建防火墙规则。
func (fw *Firewall) list(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := fw.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

// create 创建防火墙规则。
func (fw *Firewall) create(c *ship.Context) error {
	req := new(request.FirewallUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return fw.svc.Create(ctx, req)
}

// update 修改防火墙规则。
func (fw *Firewall) update(c *ship.Context) error {
	req := new(request.FirewallUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return fw.svc.Update(ctx, req)
}

// delete 删除防火墙规则。
func (fw *Firewall) delete(c *ship.Context) error {
	req := new(request.Names)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return fw.svc.Delete(ctx, req.Name)
}

// precheck 预检防火墙规则。
//
// 在管理员配置防火墙时，如果操作不当，极有可能配置完后把自己据之在外，
// 该接口就是在配置提交前，客户端将配置预检查一下，避免提交不当的配置。
// 当然该接口是可选的，在 create update 时不会检查防火墙配置是否得当，
// 此接口通常在前端提交前，自动请求该接口，如果配置不当给出友好的提示。
func (fw *Firewall) precheck(c *ship.Context) error {
	req := new(request.FirewallUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}

	if !req.Enabled {
		ret := &response.FirewallPrecheck{Allowed: true}
		return c.JSON(http.StatusOK, ret)
	}

	box, err := fw.svc.Sandbox(req)
	if err != nil {
		return err
	}
	allowed := box.Allow(c.Request())
	ret := &response.FirewallPrecheck{Allowed: allowed}

	return c.JSON(http.StatusOK, ret)
}

func (fw *Firewall) reset(c *ship.Context) error {
	fw.svc.Reset()
	return c.NoContent(http.StatusNoContent)
}
