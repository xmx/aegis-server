package gateway

import (
	"github.com/xmx/aegis-server/channel/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type brokerPeer struct {
	id   bson.ObjectID
	name string
	mux  transport.Muxer
}

func (bp *brokerPeer) ID() string             { return bp.id.Hex() }
func (bp *brokerPeer) Name() string           { return bp.name }
func (bp *brokerPeer) Muxer() transport.Muxer { return bp.mux }
