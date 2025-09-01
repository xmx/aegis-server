package broker

import (
	"github.com/xmx/aegis-server/channel/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Peer transport.Peer[bson.ObjectID]

type brokPeer struct {
	id  bson.ObjectID
	mux transport.Muxer
}

func (brk *brokPeer) ID() bson.ObjectID {
	return brk.id
}

func (brk *brokPeer) Mux() transport.Muxer {
	return brk.mux
}
