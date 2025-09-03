package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Agent struct {
	ID        bson.ObjectID `json:"id,omitzero"     bson:"_id,omitempty"`        // ID
	MachineID string        `json:"machine_id"      bson:"machine_id"`           // 机器码（全局唯一）
	Status    AgentStatus   `json:"status"          bson:"status"`               // 节点状态
	Broker    *AgentBroker  `json:"broker,omitzero" bson:"broker,omitempty"`     // agent 所在的 broker
	CreatedAt time.Time     `json:"created_at"      bson:"created_at,omitempty"` // 创建时间
	UpdatedAt time.Time     `json:"updated_at"      bson:"updated_at,omitempty"` // 更新时间
}

const (
	AgentStatusOffline AgentStatus = iota
	AgentStatusOnline
)

type AgentStatus int

func (s AgentStatus) Online() bool {
	return s == AgentStatusOnline
}

type AgentBroker struct {
	ID   bson.ObjectID `json:"id,omitzero"   bson:"id,omitempty"`
	Name string        `json:"name,omitzero" bson:"name,omitempty"`
}
