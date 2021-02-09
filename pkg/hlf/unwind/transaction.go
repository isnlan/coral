package unwind

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/isnlan/coral/pkg/protos"

	"github.com/isnlan/coral/pkg/contract/identity"

	"github.com/isnlan/coral/pkg/errors"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

type Transaction struct {
	TxId              string               `db:"tx_id" json:"tx_id"`                           // 交易ID
	ChannelId         string               `db:"channel_id" json:"channel_id"`                 // 通道ID
	BlockNumber       uint64               `db:"block_number" json:"block_number"`             // 块高度
	Timestamp         time.Time            `db:"timestamp" json:"timestamp"`                   // 时间戳
	ValidationCode    string               `db:"validation_code" json:"validation_code"`       // 交易状态
	TransactionType   string               `db:"transaction_type" json:"transaction_type"`     // 交易类型
	CreatorMspId      string               `db:"creator_msp_id" json:"creator_msp_id"`         // 交易创建者MSP ID
	CreatorIdBytes    string               `db:"creator_id_bytes" json:"creator_id_bytes"`     // 交易创建者身份ID
	CreatorNonce      []byte               `db:"creator_nonce" json:"creator_nonce"`           // 噪声
	RWSet             []byte               `db:"rwset" json:"rwset"`                           // 读写集
	Proposal          []byte               `db:"proposal" json:"proposal"`                     // 交易提案
	Response          []byte               `db:"response" json:"response"`                     // 交易结果
	ProposalHash      []byte               `db:"proposal_hash" json:"proposal_hash"`           // 提案Hash
	EnvelopeSignature []byte               `db:"envelope_signature" json:"envelope_signature"` // 交易签名
	Events            *peer.ChaincodeEvent `db:"events" json:"events"`                         // 交易事件
}

func NewTransactionFromPayload(payload []byte, validationCode int32) (*Transaction, error) {
	env, err := protoutil.GetEnvelopeFromBlock(payload)
	if err != nil {
		return nil, err
	}
	return NewTransactionFromEnvelope(env, validationCode)
}

func NewTransactionFromEnvelope(envelope *common.Envelope, validationCode int32) (*Transaction, error) {
	transaction := new(Transaction)
	payload := new(common.Payload)
	header := new(common.ChannelHeader)
	ex := &peer.ChaincodeHeaderExtension{}

	if err := proto.Unmarshal(envelope.Payload, payload); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(payload.Header.ChannelHeader, header); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(header.Extension, ex); err != nil {
		return nil, err
	}

	transaction.TxId = header.TxId
	transaction.Timestamp = time.Unix(header.Timestamp.Seconds, int64(header.Timestamp.Nanos))
	transaction.ChannelId = header.ChannelId
	transaction.ValidationCode = peer.TxValidationCode_name[validationCode]
	transaction.TransactionType = common.HeaderType_name[header.Type]
	transaction.EnvelopeSignature = envelope.Signature

	if common.HeaderType(header.Type) == common.HeaderType_ENDORSER_TRANSACTION {
		tx, err := protoutil.UnmarshalTransaction(payload.Data)
		if err != nil {
			return nil, err
		}

		if len(tx.Actions) > 0 {
			ccActionPayload, err := protoutil.UnmarshalChaincodeActionPayload(tx.Actions[0].Payload)
			if err != nil {
				return nil, err
			}

			signatureHeader, err := protoutil.UnmarshalSignatureHeader(tx.Actions[0].Header)
			if err != nil {
				return nil, err
			}

			identity := &msp.SerializedIdentity{}
			err = proto.Unmarshal(signatureHeader.Creator, identity)
			if err != nil {
				return nil, err
			}

			transaction.CreatorMspId = identity.Mspid
			transaction.CreatorIdBytes = string(identity.IdBytes)
			transaction.CreatorNonce = signatureHeader.Nonce

			ccProposalPayload, err := protoutil.UnmarshalChaincodeProposalPayload(ccActionPayload.ChaincodeProposalPayload)
			if err != nil {
				return nil, err
			}
			transaction.Proposal = ccProposalPayload.Input

			propRespPayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(ccActionPayload.Action.ProposalResponsePayload, propRespPayload)
			if err != nil {
				return nil, err
			}

			ccAction := &peer.ChaincodeAction{}
			err = proto.Unmarshal(propRespPayload.Extension, ccAction)
			if err != nil {
				return nil, err
			}

			transaction.RWSet = ccAction.Results
			transaction.ProposalHash = propRespPayload.ProposalHash

			if ccAction.Response == nil ||
				ccAction.Response.Payload == nil ||
				len(ccAction.Response.Payload) == 0 {
				transaction.Response = []byte{}
			}

			if len(ccAction.Events) != 0 {
				ccEvent, err := protoutil.UnmarshalChaincodeEvents(ccAction.Events)
				if err != nil {
					return nil, errors.WithMessage(err, "error unmarshal chaincode event for block event")
				}
				transaction.Events = ccEvent
			}
		}
	}

	return transaction, nil
}

func (t *Transaction) IntoTransaction() *protos.InnerTransaction {
	timestamp, _ := ptypes.TimestampProto(t.Timestamp)

	tx := protos.InnerTransaction{
		TxId:           t.TxId,
		ChannelId:      t.ChannelId,
		BlockNumber:    t.BlockNumber,
		Timestamp:      timestamp,
		ValidationCode: t.ValidationCode,
		Event:          nil,
	}
	if t.Events != nil {
		e := protos.Event{
			Contract:  t.Events.ChaincodeId,
			EventName: t.Events.EventName,
			Value:     t.Events.Payload,
		}
		tx.Event = &e
	}

	if proposal, err := protoutil.UnmarshalChaincodeInvocationSpec(t.Proposal); err == nil {
		if proposal.ChaincodeSpec != nil && proposal.ChaincodeSpec.ChaincodeId != nil {
			tx.Contract = proposal.ChaincodeSpec.ChaincodeId.Name
		}
	}

	tx.Sign = t.EnvelopeSignature
	tx.TxType = t.TransactionType

	if address, err := identity.IntoAddress([]byte(t.CreatorIdBytes)); err == nil {
		tx.Creator = address.String()
	}
	return &tx
}
