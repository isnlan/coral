package unwind

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/snlansky/coral/pkg/protos"

	"github.com/hyperledger/fabric-protos-go/common"
)

func NewBlock(block *common.Block) (*protos.Block, error) {
	tb := new(protos.Block)
	tb.Number = block.Header.Number
	tb.PreviousHash = block.Header.PreviousHash
	tb.Hash = protoutil.BlockHeaderHash(block.Header)
	tb.DataHash = block.Header.DataHash

	bytes, err := protoutil.Marshal(block)
	if err != nil {
		return nil, err
	}
	tb.Size = int32(len(bytes))

	var txs []*protos.Transaction

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
			return nil, err
		}

		tb.ChannelId = transaction.ChannelId
		tb.Timestamp, err = ptypes.TimestampProto(transaction.Timestamp)
		if err != nil {
			return nil, err
		}

		// 只保存验证通过的
		if transaction.ValidationCode != peer.TxValidationCode_name[0] {
			continue
		}

		// 只保存交易，其他的全部虐过
		if transaction.TransactionType != common.HeaderType_name[3] {
			continue
		}

		transaction.BlockNumber = tb.Number
		txs = append(txs, transaction.IntoTransaction())
	}

	tb.Transactions = txs

	return tb, nil
}
