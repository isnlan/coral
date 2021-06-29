package entity

import (
	"github.com/isnlan/coral/pkg/protos"
)

type Block struct {
	Number           uint64 `db:"number" json:"number" bson:"number"`                                  // 当前块高度
	PreviousHash     []byte `db:"previous_hash" json:"previous_hash" bson:"previous_hash"`             // 前块Hash
	Hash             []byte `db:"hash" json:"hash" bson:"hash"`                                        // 当前块Hash
	DataHash         []byte `db:"data_hash" json:"data_hash" bson:"data_hash"`                         // 数据Hash
	TransactionCount int    `db:"transaction_count" json:"transaction_count" bson:"transaction_count"` // 当前块中交易数量
	ChannelId        string `db:"channel_id" json:"channel_id" bson:"channel_id"`                      // 链ID
	Size             int    `db:"size" json:"size" bson:"size"`                                        // 区块大小
	Timestamp        int64  `db:"timestamp" json:"timestamp" bson:"timestamp"`                         // 时间戳
}

//----------------------------------------------------------------------------------------------------------------------
type Transaction struct {
	TxId           string `db:"tx_id" json:"tx_id" bson:"tx_id"`                               // 交易ID
	ChannelId      string `db:"channel_id" json:"channel_id" bson:"channel_id"`                // 通道ID
	BlockNumber    uint64 `db:"block_number" json:"block_number" bson:"block_number"`          // 块高度
	Contract       string `db:"contract" json:"contract" bson:"contract"`                      // 合约名称
	Creator        string `db:"creator" json:"creator" bson:"creator"`                         // 数字身份
	Sign           []byte `db:"sign" json:"sign" bson:"sign"`                                  // 签名
	TxType         string `db:"tx_type" json:"tx_type" bson:"tx_type" `                        // 交易类型
	Timestamp      int64  `db:"timestamp" json:"timestamp" bson:"timestamp"`                   // 时间戳
	ValidationCode string `db:"validation_code" json:"validation_code" bson:"validation_code"` // 交易状态
	Event          *Event `db:"event" json:"event" bson:"event"`                               // 事件
}

//----------------------------------------------------------------------------------------------------------------------
type Event struct {
	Contract  string `json:"contract" bson:"contract"`
	EventName string `json:"event_name" bson:"event_name"`
	Value     []byte `json:"value" bson:"value"`
}

func FromInnerBlock(b *protos.InnerBlock) (*Block, []*Transaction) {
	block := &Block{
		Number:           b.Number,
		PreviousHash:     b.PreviousHash,
		Hash:             b.Hash,
		DataHash:         b.DataHash,
		TransactionCount: len(b.Transactions),
		ChannelId:        b.ChannelId,
		Size:             int(b.Size),
		Timestamp:        b.Timestamp.AsTime().Unix(),
	}

	var txs []*Transaction
	for _, tx := range b.Transactions {
		txs = append(txs, FromInnerTransaction(tx))
	}

	return block, txs
}

func FromInnerTransaction(t *protos.InnerTransaction) *Transaction {
	var event *Event
	if t.Event != nil {
		event = &Event{
			Contract:  t.Event.Contract,
			EventName: t.Event.EventName,
			Value:     t.Event.Value,
		}
	}
	return &Transaction{
		TxId:           t.TxId,
		ChannelId:      t.ChannelId,
		BlockNumber:    t.BlockNumber,
		Contract:       t.Contract,
		Creator:        t.Creator,
		Sign:           t.Sign,
		TxType:         t.TxType,
		Timestamp:      t.Timestamp.AsTime().Unix(),
		ValidationCode: t.ValidationCode,
		Event:          event,
	}
}
