package model_utils

import (

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/snlansky/coral/pkg/model"

	"github.com/hyperledger/fabric-protos-go/common"
)

func NewBlock(block *common.Block) (*model.Block, []*Transaction, error) {
	tb := new(model.Block)
	tb.Number = block.Header.Number
	tb.PreviousHash = block.Header.PreviousHash
	tb.Hash = protoutil.BlockHeaderHash(block.Header)
	tb.DataHash = block.Header.DataHash

	var txs []*Transaction

	for idx, pl := range block.Data.Data {
		var status byte
		if len(block.Metadata.Metadata) > 2 {
			meta := block.Metadata.Metadata[2]
			if len(meta) > idx {
				status = meta[idx]
			}
		}

		transaction, err := NewTransactionFromPayload(pl, int32(status))
		if err != nil {
			return nil, nil, err
		}

		tb.ChannelId = transaction.ChannelId
		tb.Timestamp = transaction.Timestamp

		// 只保存验证通过的
		if transaction.ValidationCode != peer.TxValidationCode_name[0] {
			continue
		}

		// 只保存交易，其他的全部虐过
		if transaction.TransactionType != common.HeaderType_name[3] {
			continue
		}

		transaction.BlockNumber = tb.Number
		txs = append(txs, transaction)
	}

	tb.TransactionCount = len(txs)
	return tb, txs, nil
}
