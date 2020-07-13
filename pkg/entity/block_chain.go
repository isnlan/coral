package entity

import (
	"time"
)

type Block struct {
	Number           uint64    `db:"number" json:"number" bson:"number"`                                  // 当前块高度
	PreviousHash     []byte    `db:"previous_hash" json:"previous_hash" bson:"previous_hash"`             // 前块Hash
	Hash             []byte    `db:"hash" json:"hash" bson:"hash"`                                        // 当前块Hash
	DataHash         []byte    `db:"data_hash" json:"data_hash" bson:"data_hash"`                         // 数据Hash
	TransactionCount int       `db:"transaction_count" json:"transaction_count" bson:"transaction_count"` // 当前块中交易数量
	ChannelId        string    `db:"channel_id" json:"channel_id" bson:"channel_id"`                      // 链ID
	Timestamp        time.Time `db:"timestamp" json:"timestamp" bson:"timestamp"`                         // 时间戳
}

//----------------------------------------------------------------------------------------------------------------------
type Transaction struct {
	TxId           string    `db:"tx_id" json:"tx_id" bson:"tx_id"`                               // 交易ID
	ChannelId      string    `db:"channel_id" json:"channel_id" bson:"channel_id"`                // 通道ID
	BlockNumber    uint64    `db:"block_number" json:"block_number" bson:"block_number"`          // 块高度
	Contract       string    `db:"contract" json:"contract" bson:"contract"`                      // 合约名称
	Creator        string    `db:"creator" json:"creator" bson:"creator"`                         // 数字身份
	Sign           []byte    `db:"sign" json:"sign" bson:"sign"`                                  // 签名
	TxType         string    `db:"tx_type" json:"tx_type" bson:"tx_type" `                        // 交易类型
	Timestamp      time.Time `db:"timestamp" json:"timestamp" bson:"timestamp"`                   // 时间戳
	ValidationCode string    `db:"validation_code" json:"validation_code" bson:"validation_code"` // 交易状态
	Event          *Event    `db:"event" json:"event" bson:"event"`                               // 事件
}

//----------------------------------------------------------------------------------------------------------------------
type Event struct {
	Contract  string `json:"contract" bson:"contract"`
	EventName string `json:"event_name" bson:"event_name"`
	Value     []byte `json:"value" bson:"value"`
}
