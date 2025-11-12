package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewBrokerRelease(svc *service.BrokerRelease) *BrokerRelease {
	return &BrokerRelease{
		svc: svc,
	}
}

type BrokerRelease struct {
	svc *service.BrokerRelease
}

func (br *BrokerRelease) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/broker-release").
		GET(br.download).
		POST(br.upload)
	r.Route("/broker-release/parse").POST(br.parse)

	return nil
}

func (br *BrokerRelease) upload(c *ship.Context) error {
	req := new(request.BrokerReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return br.svc.Upload(ctx, req)
}

func (br *BrokerRelease) parse(c *ship.Context) error {
	req := new(request.BrokerReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, ret, err := br.svc.Parse(file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (br *BrokerRelease) download(c *ship.Context) error {
	//req := new(request.ObjectID)
	//if err := c.BindQuery(req); err != nil {
	//	return err
	//}
	//
	//ctx := c.Request().Context()
	//info, stm, err := br.svc.Open(ctx, req.OID())
	//if err != nil {
	//	return err
	//}
	//defer stm.Close()
	//
	//exposes, err := br.brok.Exposes(ctx)
	//if err != nil {
	//	return err
	//}
	//
	//filesize := info.Length
	//manifest := &response.AgentManifest{
	//	Addresses: exposes,
	//	Offset:    filesize,
	//}
	//
	//zipbuf, err := stegano.CreateManifestZip(manifest, filesize)
	//if err != nil {
	//	return err
	//}
	//
	//totalLen := filesize + int64(zipbuf.Len())
	//contentLength := strconv.FormatInt(totalLen, 10)
	//params := info.Checksum.Map()
	//params["filename"] = info.Filename
	//mediaType := mime.FormatMediaType("attachment", params)
	//c.SetRespHeader(ship.HeaderContentDisposition, mediaType)
	//c.SetRespHeader(ship.HeaderContentLength, contentLength)
	//down := io.MultiReader(stm, zipbuf)
	//
	//return c.Stream(http.StatusOK, ship.MIMEOctetStream, down)
	return nil
}
