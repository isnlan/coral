/*
Copyright: Cognition Foundry. All Rights Reserved.
License: Apache License Version 2.0
*/
package hlf

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// Peer expose API's to communicate with peer
type Peer struct {
	Name               string
	Uri                string
	MspId              string
	Opts               []grpc.DialOption
	Cert               string
	ServerNameOverride string
	conn               *grpc.ClientConn
	client             peer.EndorserClient
}

// PeerResponse is response from peer transaction request
type PeerResponse struct {
	Response *peer.ProposalResponse
	Err      error
	Name     string
}

// Endorse sends single transaction to single peer.
func (p *Peer) Endorse(ctx context.Context, resp chan *PeerResponse, prop *peer.SignedProposal) {
	if p.conn == nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		conn, err := grpc.DialContext(ctx, p.Uri, p.Opts...)
		if err != nil {
			resp <- &PeerResponse{Response: nil, Err: err, Name: p.Name}
			return
		}
		p.conn = conn
		p.client = peer.NewEndorserClient(p.conn)
	}
	proposalResp, err := p.client.ProcessProposal(ctx, prop)
	if err != nil {
		resp <- &PeerResponse{Response: nil, Name: p.Name, Err: err}
		return
	}
	if proposalResp.Response.Status != 200 {
		resp <- &PeerResponse{Response: nil, Name: p.Name, Err: errors.New(proposalResp.Response.Message)}
		return
	}
	resp <- &PeerResponse{Response: proposalResp, Name: p.Name, Err: nil}
}

// NewPeerFromConfig creates new peer from provided config
func NewPeerFromConfig(conf PeerConfig) (*Peer, error) {
	p := Peer{Uri: conf.Host, Cert: conf.Cert, ServerNameOverride: conf.ServerNameOverride}
	if !conf.UseTLS {
		p.Opts = []grpc.DialOption{grpc.WithInsecure()}
	} else if p.Cert != "" {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM([]byte(p.Cert)) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
		creds := credentials.NewClientTLSFromCert(cp, p.ServerNameOverride)
		p.Opts = append(p.Opts, grpc.WithTransportCredentials(creds))
	}

	p.Opts = append(p.Opts,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(1) * time.Minute,
			Timeout:             time.Duration(20) * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize)))
	return &p, nil
}
