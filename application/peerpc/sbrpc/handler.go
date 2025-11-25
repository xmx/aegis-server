package sbrpc

import (
	"context"
	"net/url"

	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/tunnel/tunconst"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Handler interface {
	SystemExit(ctx context.Context) error
}

type handler struct {
	cli httpkit.Client
	bid bson.ObjectID
}

func (h *handler) SystemExit(ctx context.Context) error {
	reqURL := h.buildURL("/system/exit")
	strURL := reqURL.String()

	return h.cli.GetJSON(ctx, strURL, nil, nil)
}

func (h *handler) buildURL(path string) *url.URL {
	bid := h.bid.Hex()
	u := tunconst.ServerToBroker(bid, "/api/")
	return u.JoinPath(path)
}
