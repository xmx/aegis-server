package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Agent struct {
	ID        bson.ObjectID `json:"id,omitzero"           bson:"_id,omitempty"` // ID
	MachineID string        `json:"machine_id"            bson:"machine_id"`    // 机器码（全局唯一）
	Status    bool          `json:"status"                bson:"status"`        // 节点状态
	//Networks    NodeNetworks  `json:"networks,omitzero"     bson:"networks,omitempty"`     // 网卡
	//TunnelStat  *TunnelStat   `json:"tunnel_stat,omitzero"  bson:"tunnel_stat,omitempty"`  // 通道连接状态
	//ExecuteStat *ExecuteStat  `json:"execute_stat,omitzero" bson:"execute_stat,omitempty"` // agent 执行程序信息
	CreatedAt time.Time `json:"created_at"            bson:"created_at,omitempty"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"            bson:"updated_at,omitempty"` // 更新时间
}

type Agents []*Agent
