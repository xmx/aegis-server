package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AgentConnectHistory struct {
	ID         bson.ObjectID     `json:"id"              bson:"_id,omitempty"`
	AgentID    bson.ObjectID     `json:"agent_id"        bson:"agent_id"`
	MachineID  string            `json:"machine_id"      bson:"machine_id"`
	Semver     string            `json:"semver,omitzero" bson:"semver,omitempty"`
	Inet       string            `json:"inet,omitzero"   bson:"inet,omitempty"`
	Goos       string            `json:"goos,omitzero"   bson:"goos,omitempty"`
	Goarch     string            `json:"goarch,omitzero" bson:"goarch,omitempty"`
	TunnelStat TunnelStatHistory `json:"tunnel_stat"     bson:"tunnel_stat"`
}

type TunnelStatHistory struct {
	ConnectedAt    time.Time     `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`
	DisconnectedAt time.Time     `json:"disconnected_at,omitzero" bson:"disconnected_at,omitempty"`
	Second         int64         `json:"second"                   bson:"second"`
	Library        TunnelLibrary `json:"library"                  bson:"library"`
	LocalAddr      string        `json:"local_addr,omitzero"      bson:"local_addr,omitempty"`
	RemoteAddr     string        `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`
	ReceiveBytes   uint64        `json:"receive_bytes,omitzero"   bson:"receive_bytes,omitempty"`  // broker/agent 为主体
	TransmitBytes  uint64        `json:"transmit_bytes,omitzero"  bson:"transmit_bytes,omitempty"` // broker/agent 为主体
}

type TunnelLibrary struct {
}
