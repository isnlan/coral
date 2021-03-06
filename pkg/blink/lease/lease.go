package lease

type ChainLease struct {
	NetworkID   string `json:"network_id" mapstructure:"network_id"`
	NetworkType string `json:"network_type" mapstructure:"network_type"`
	NetworkName string `json:"network_name" mapstructure:"network_name"`
	Account     string `json:"account" mapstructure:"account"`
	Team        string `json:"team" mapstructure:"team"`
	IsRunning   bool   `json:"is_running" mapstructure:"is_running"`
	Status      string `json:"status" mapstructure:"status"`
	TlsEnabled  bool   `json:"tls_enabled" mapstructure:"tls_enabled"`
}

func (l *ChainLease) UniqueID() string {
	return l.NetworkID
}

type ChannelLease struct {
	ID         string `json:"id" mapstructure:"id"`
	NetworkID  string `json:"network_id" mapstructure:"network_id"`
	Name       string `json:"name" mapstructure:"name"`
	Endpoint   string `json:"endpoint" mapstructure:"endpoint"`
	IsRunning  bool   `json:"is_running" mapstructure:"is_running"`
	SyncEnable bool   `json:"sync_enable" mapstructure:"sync_enable"`
	SyncDB     string `json:"sync_db" mapstructure:"sync_db"`
}

func (l *ChannelLease) UniqueID() string {
	return l.NetworkID + ":" + l.Name
}

type AclLease struct {
	ID           string   `json:"id" mapstructure:"id"`                       // 应用ID
	Name         string   `json:"name" mapstructure:"name"`                   // 应用名称
	ClientId     string   `json:"client_id" mapstructure:"client_id"`         // 客户端ID
	ClientSecret string   `json:"client_secret" mapstructure:"client_secret"` // 客户端Secret
	Account      string   `json:"account" mapstructure:"account"`             // 创建账户
	Team         string   `json:"team" mapstructure:"team"`                   // 组
	NetworkID    string   `json:"network_id" mapstructure:"network_id"`       // 网络ID（ChainID）
	Nodes        []string `json:"nodes" mapstructure:"nodes"`                 // 可用节点
	CreateTime   int64    `json:"create_time" mapstructure:"mapstructure"`    // 创建时间
	Enable       bool     `json:"enable" mapstructure:"enable"`               // 可用
	Description  string   `json:"description" mapstructure:"description"`     // 描述
}

func (l *AclLease) UniqueID() string {
	return l.ClientId
}
