package sbrpc

import (
	"github.com/xmx/aegis-common/library/httpkit"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewClient(cli httpkit.Client) Client {
	return Client{cli: cli}
}

type Client struct {
	cli httpkit.Client
}

func (c Client) NewHandler(brokID bson.ObjectID) Handler {
	return &handler{
		cli: c.cli, bid: brokID,
	}
}
