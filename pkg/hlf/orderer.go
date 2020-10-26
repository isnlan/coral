/*
Copyright: Cognition Foundry. All Rights Reserved.
License: Apache License Version 2.0
*/
package hlf

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/snlansky/coral/pkg/trace"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// Orderer expose API's to communicate with orderers.
type Orderer struct {
	Name               string
	Uri                string
	Opts               []grpc.DialOption
	Cert               string
	ServerNameOverride string
	conn               *grpc.ClientConn
	client             orderer.AtomicBroadcastClient
}

const timeout = 5

// Broadcast Broadcast envelope to orderer for execution.
func (o *Orderer) Broadcast(envelope *common.Envelope) (*orderer.BroadcastResponse, error) {
	if o.conn == nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		c, err := grpc.DialContext(ctx, o.Uri, o.Opts...)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to orderer: %s err is: %v", o.Name, err)
		}
		o.conn = c
		o.client = orderer.NewAtomicBroadcastClient(o.conn)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	bcc, err := o.client.Broadcast(ctx)
	if err != nil {
		return nil, err
	}
	defer bcc.CloseSend()
	bcc.Send(envelope)
	response, err := bcc.Recv()
	if err != nil {
		return nil, err
	}
	if response.Status != common.Status_SUCCESS {
		return nil, fmt.Errorf("unexpected status: %v", response.Status)
	}

	return response, err
}

// Deliver delivers envelope to orderer. Please note that new connection will be created on every call of Deliver.
func (o *Orderer) Deliver(envelope *common.Envelope) (*common.Block, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	connection, err := grpc.DialContext(ctx, o.Uri, o.Opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to orderer: %s err is: %v", o.Name, err)
	}
	defer connection.Close()

	dk, err := orderer.NewAtomicBroadcastClient(connection).Deliver(context.Background())
	if err != nil {
		return nil, err
	}
	if err := dk.Send(envelope); err != nil {
		return nil, err
	}
	var block *common.Block
	timer := time.NewTimer(time.Second * time.Duration(timeout))
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return nil, ErrOrdererTimeout
		default:
			response, err := dk.Recv()
			if err != nil {
				return nil, err
			}
			switch t := response.Type.(type) {
			case *orderer.DeliverResponse_Status:
				if t.Status == common.Status_SUCCESS {
					return block, nil
				} else {
					return nil, fmt.Errorf("orderer response with status: %v", t.Status)
				}
			case *orderer.DeliverResponse_Block:
				block = response.GetBlock()

			default:
				return nil, fmt.Errorf("unknown response type from orderer: %s", t)
			}
		}
	}
}

func (o *Orderer) getGenesisBlock(identity Identity, crypto CryptoSuite, channelId string) (*common.Block, error) {

	seekInfo := &orderer.SeekInfo{
		Start:    &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: 0}}},
		Stop:     &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: 0}}},
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}
	seekInfoBytes, err := proto.Marshal(seekInfo)
	if err != nil {
		return nil, err
	}

	creator, err := marshalProtoIdentity(identity)
	if err != nil {
		return nil, err
	}
	txId, err := newTransactionId(creator)
	if err != nil {
		return nil, err
	}

	headerBytes, err := channelHeader(common.HeaderType_DELIVER_SEEK_INFO, txId, channelId, 0, nil)
	signatureHeaderBytes, err := signatureHeader(creator, txId)
	if err != nil {
		return nil, err
	}
	header := header(signatureHeaderBytes, headerBytes)
	payloadBytes, err := payload(header, seekInfoBytes)
	if err != nil {
		return nil, err
	}
	payloadSignedBytes, err := crypto.Sign(payloadBytes, identity.PrivateKey)
	if err != nil {
		return nil, err
	}
	env := &common.Envelope{Payload: payloadBytes, Signature: payloadSignedBytes}
	return o.Deliver(env)
}

// NewOrdererFromConfig create new Orderer from config
func NewOrdererFromConfig(conf OrdererConfig) (*Orderer, error) {
	o := Orderer{Uri: conf.Host, Cert: conf.Cert, ServerNameOverride: conf.ServerNameOverride}
	if !conf.UseTLS {
		o.Opts = []grpc.DialOption{grpc.WithInsecure()}
	} else if o.Cert != "" {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM([]byte(o.Cert)) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
		creds := credentials.NewClientTLSFromCert(cp, o.ServerNameOverride)
		o.Opts = append(o.Opts, grpc.WithTransportCredentials(creds))
	}
	o.Opts = append(o.Opts,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(1) * time.Minute,
			Timeout:             time.Duration(20) * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize)),
		grpc.WithUnaryInterceptor(
			grpc_middeware.ChainUnaryClient(
				trace.OpenTracingClientInterceptor(),
			),
		),
	)
	return &o, nil
}
