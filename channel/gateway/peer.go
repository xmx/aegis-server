package gateway

import (
	"github.com/xmx/aegis-server/channel/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type brokerPeer struct {
	id     bson.ObjectID
	mux    transport.Muxer
	goos   string
	goarch string
}

func (bp *brokerPeer) ID() string             { return bp.id.Hex() }
func (bp *brokerPeer) Muxer() transport.Muxer { return bp.mux }
func (bp *brokerPeer) Goos() string           { return bp.goos }
func (bp *brokerPeer) Goarch() string         { return bp.goarch }
