package model

import "time"

type CheckPoint struct {
	NetworkId     string    `json:"network_id" bson:"network_id"`
	ChainId       string    `json:"chain_id" bson:"chain_id"`
	SyncedBlock   uint64    `json:"synced_block" bson:"synced_block"`
	SyncedTotalTx uint64    `json:"synced_total_tx" bson:"synced_total_tx"`
	SyncedTime    time.Time `json:"synced_time" bson:"synced_time"`
}
