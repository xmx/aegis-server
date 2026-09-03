package model

import "go.mongodb.org/mongo-driver/v2/bson"

type AgentConfig struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}

func (AgentConfig) CollectionInfo() (string, string) {
	return "agent_config", "Agent 配置"
}
