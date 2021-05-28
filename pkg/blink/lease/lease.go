package lease

type ChainLease struct {
	NetworkID   string `json:"network_id"`
	NetworkType string `json:"network_type"`
	NetworkName string `json:"network_name"`
	Account     string `json:"account"`
	Team        string `json:"team"`
	IsRunning   bool   `json:"is_running"`
	Status      string `json:"status"`
	TlsEnabled  bool   `json:"tls_enabled"`
}

type ChannelLease struct {
	ID         string `json:"id"`
	NetworkID  string `json:"network_id"`
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	IsRunning  bool   `json:"status"`
	SyncEnable bool   `json:"sync_enable"`
	SyncDB     string `json:"sync_db"`
}

type AclLease struct {
	ID           string `json:"id"`            // 应用ID
	Name         string `json:"name"`          // 应用名称
	ClientId     string `json:"client_id"`     // 客户端ID
	ClientSecret string `json:"client_secret"` // 客户端Secret
	Account      string `json:"account"`       // 创建账户
	Team         string `json:"team"`          // 组
	NetworkID    string `json:"network_id"`    // 网络ID（ChainID）
	Enable       bool   `json:"enable"`        // 可用
}
