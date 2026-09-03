package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Agent struct {
	ID          bson.ObjectID    `json:"id,omitzero"  bson:"_id,omitempty"` // ID
	MachineID   string           `json:"machine_id"   bson:"machine_id"`    // 机器码（全局唯一）
	Status      bool             `json:"status"       bson:"status"`        // 节点状态
	Enabled     bool             `json:"enabled"      bson:"enabled"`       // 是否启用（禁用的节点不允许上线）
	ProcessInfo AgentProcessInfo `json:"process_info" bson:"process_info"`
	TunnelInfo  AgentTunnelInfo  `json:"tunnel_info"  bson:"tunnel_info"`
	CreatedAt   time.Time        `json:"created_at"   bson:"created_at,omitempty"` // 创建时间
	UpdatedAt   time.Time        `json:"updated_at"   bson:"updated_at,omitempty"` // 更新时间
}

func (Agent) CollectionInfo() (string, string) { return "agent", "Agent 节点列表" }

type AgentProcessInfo struct {
	Inet       string   `json:"inet"                bson:"inet"`
	Semver     Semver   `json:"semver,omitzero"     bson:"semver,omitempty"`
	GOOS       string   `json:"goos,omitzero"       bson:"goos,omitempty"`
	GOARCH     string   `json:"goarch,omitzero"     bson:"goarch,omitempty"`
	PID        int      `json:"pid,omitzero"        bson:"pid,omitempty"`
	Args       []string `json:"args,omitzero"       bson:"args,omitempty"`
	Hostname   string   `json:"hostname,omitzero"   bson:"hostname,omitempty"`
	Workdir    string   `json:"workdir,omitzero"    bson:"workdir,omitempty"`
	Executable string   `json:"executable,omitzero" bson:"executable,omitempty"`
	Environ    []string `json:"environ,omitzero"    bson:"environ,omitempty"`
}

type AgentTunnelInfo struct {
	ConnectedAt   time.Time `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`
	KeepaliveAt   time.Time `json:"keepalive_at,omitzero"    bson:"keepalive_at,omitempty"`
	LibraryName   string    `json:"library_name"             bson:"library_name"`
	LibraryModule string    `json:"library_module"           bson:"library_module"`
	ServerAddr    string    `json:"server_addr"              bson:"server_addr"`
	RemoteAddr    string    `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`
	ReceiveBytes  uint64    `json:"receive_bytes"            bson:"receive_bytes,omitempty"`  // broker/agent 为主体
	TransmitBytes uint64    `json:"transmit_bytes"           bson:"transmit_bytes,omitempty"` // broker/agent 为主体
}
