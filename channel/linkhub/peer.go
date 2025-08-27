package linkhub

import "github.com/xmx/aegis-server/channel/transport"

type Peer interface {
	// ID 节点全局唯一 ID。
	ID() string

	// Muxer 底层多路通道。
	Muxer() transport.Muxer
}

type brokerPeer struct {
	id  string
	mux transport.Muxer
}

func (bp *brokerPeer) ID() string {
	return bp.id
}

func (bp *brokerPeer) Muxer() transport.Muxer {
	return bp.mux
}
