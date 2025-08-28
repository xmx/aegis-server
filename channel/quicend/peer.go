package quicend

import "github.com/xmx/aegis-server/channel/transport"

type agentPeer struct {
	id  string
	mux transport.Muxer
}

func (a *agentPeer) ID() string {
	return a.id
}

func (a *agentPeer) Muxer() transport.Muxer {
	return a.mux
}
