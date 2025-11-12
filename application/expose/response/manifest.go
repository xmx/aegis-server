package response

type AgentManifest struct {
	Protocols []string `json:"protocols"` // 连接协议 udp tcp，一般留空即可。
	Addresses []string `json:"addresses"` // 连接的 broker 地址
	Offset    int64    `json:"offset"`
}

type BrokerManifest struct {
	Protocols []string `json:"protocols"`
	Addresses []string `json:"addresses"`
	Secret    string   `json:"secret"`
}
