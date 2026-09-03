package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AgentConnRecord struct {
	ID             bson.ObjectID    `json:"id"              bson:"_id,omitempty"`
	AgentID        bson.ObjectID    `json:"agent_id"        bson:"agent_id"`
	MachineID      string           `json:"machine_id"      bson:"machine_id"`
	ProcessInfo    AgentProcessInfo `json:"process_info"    bson:"process_info"`
	TunnelInfo     AgentTunnelInfo  `json:"tunnel_info"     bson:"tunnel_info"`
	ActiveSeconds  int64            `json:"active_seconds"  bson:"active_seconds"`
	DisconnectedAt time.Time        `json:"disconnected_at" bson:"disconnected_at,omitempty"` // 创建时间
}

func (AgentConnRecord) CollectionInfo() (string, string) {
	return "agent_conn_record", "Agent 连接记录"
}
