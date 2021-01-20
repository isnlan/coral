package entity

import (
	"github.com/snlansky/coral/pkg/hlf"
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

	ChannelConfig struct {
		Peers      []string `json:"peers" yaml:"peers"`
		Orderers   []string `json:"orderers" yaml:"orderers"`
		Chaincodes []string `json:"chaincodes" yaml:"chaincodes"`
		Enable     bool     `json:"enable" yaml:"enable"`
	}

	Enrollment struct {
		Cert       string `json:"cert" yaml:"cert"`
		PrivateKey string `json:"privateKey" yaml:"privateKey"`
	}

	FabricUser struct {
		Name   string     `json:"name" yaml:"name"`
		MspId  string     `json:"mspId" yaml:"mspId"`
		Enroll Enrollment `json:"enrollment" yaml:"enrollment"`
	}

	FabricConfig struct {
		Ca       *CaConfig                    `json:"ca" yaml:"ca"`
		Peers    map[string]hlf.PeerConfig    `json:"peers" yaml:"peers"`
		Orderers map[string]hlf.OrdererConfig `json:"orderers" yaml:"orderers"`
		Channels map[string]ChannelConfig     `json:"channels" yaml:"channels"`
		Admin    *FabricUser                  `json:"admin" yaml:"admin"`
	}
)

type ChannelStatus struct {
	Name   string `json:"name" yaml:"name"`
	Enable bool   `json:"enable" yaml:"enable"`
}

type Lease struct {
	NetworkID   string                    `json:"network_id"`
	NetworkType string                    `json:"network_type"`
	Account     string                    `json:"account"`
	Team        string                    `json:"team"`
	Enable      bool                      `json:"enable"`
	Channels    map[string]*ChannelStatus `json:"channels"`
	ExpireTime  int64                     `json:"expire_time"`
}
