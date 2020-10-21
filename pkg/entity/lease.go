package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AclClient struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`                      // 应用ID
	Name         string             `json:"name" bson:"name"`                   // 应用名称
	ClientId     string             `json:"client_id" bson:"client_id"`         // 客户端ID
	ClientSecret string             `json:"client_secret" bson:"client_secret"` // 客户端Secret
	Account      string             `json:"account" bson:"account"`             // 创建账户
	Team         string             `json:"team" bson:"team"`                   // 组
	ChainId      string             `json:"chain_id" bson:"chain_id"`           // ChainID
	Nodes        []string           `json:"nodes" bson:"nodes"`                 // 可用节点
	Enable       bool               `json:"enable" bson:"enable"`               // 可用
	CreateTime   int64              `json:"create_time" bson:"create_time"`     // 创建时间
	Description  string             `json:"description" bson:"description"`     // 描述
}

type (
	CaConfig struct {
		Name          string `json:"name"`
		Url           string `json:"url"`
		MspId         string `json:"mspId"`
		Affiliation   string `json:"affiliation"`
		AdminUsername string `json:"adminUsername"`
		AdminPassword string `json:"adminPassword"`
	}

	PeerConfig struct {
		Host               string `json:"host"`
		UseTLS             bool   `json:"useTLS"`
		Cert               string `json:"cert"`
		ServerNameOverride string `json:"serverNameOverride"`
	}

	OrdererConfig struct {
		Host               string `json:"host"`
		UseTLS             bool   `json:"useTLS"`
		Cert               string `json:"cert"`
		ServerNameOverride string `json:"serverNameOverride"`
	}

	ChannelConfig struct {
		Peers      []string `json:"peers"`
		Orderers   []string `json:"orderers"`
		Chaincodes []string `json:"chaincodes"`
		Enable     bool     `json:"enable"`
	}

	Enrollment struct {
		Cert       string `json:"cert"`
		PrivateKey string `json:"privateKey"`
	}

	FabricUser struct {
		Name   string     `json:"name"`
		MspId  string     `json:"mspId"`
		Enroll Enrollment `json:"enrollment"`
	}

	FabricConfig struct {
		Ca       *CaConfig                 `json:"ca"`
		Peers    map[string]*PeerConfig    `json:"peers"`
		Orderers map[string]*OrdererConfig `json:"orderers"`
		Channels map[string]*ChannelConfig `json:"channels"`
		Admin    *FabricUser               `json:"admin"`
	}
)

type Lease struct {
	NetworkId  string      `json:"network_id"`
	Consensus  string      `json:"consensus"`
	NodeCount  int         `json:"node_count"`
	Account    string      `json:"account"`
	Team       string      `json:"team"`
	TlsEnabled bool        `json:"tls_enabled"`
	Enable     bool        `json:"enable"`
	Config     interface{} `json:"config"`
}
