package main

import (
	"context"
	"fmt"

	"github.com/snlansky/coral/pkg/hlf"

	"github.com/hyperledger/fabric-protos-go/common"
)

const Cert = `-----BEGIN CERTIFICATE-----
MIICBTCCAaugAwIBAgIQZLApmNcajgGjB5iC8t9utTAKBggqhkjOPQQDAjBbMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzENMAsGA1UEChMEb3JnMTEQMA4GA1UEAxMHY2Eub3JnMTAeFw0yMDA2
MjIwMDI1MDBaFw0zMDA2MjAwMDI1MDBaMF8xCzAJBgNVBAYTAlVTMRMwEQYDVQQI
EwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ4wDAYDVQQLEwVh
ZG1pbjETMBEGA1UEAwwKQWRtaW5Ab3JnMTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABHKkQGCbDQdw8zmqyBGKpygbtyU75xPsy6osyIVTBWYsYRH4VRdkKxoqJxgZ
YW67X8T5v/81RuauoqeajJLQhvKjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMB
Af8EAjAAMCsGA1UdIwQkMCKAIPi3Pb5wquira1i886vqLO17BxyEJZASbouE9qy+
7uFsMAoGCCqGSM49BAMCA0gAMEUCIQCxe3IeRV5buyIQdqJ9S8+eMU0qwHbC93RH
+ij1M6BesAIgcykqYLxEOGGsCXDDqu/S3dfLg4SSl4JbHmfwxtwuRos=
-----END CERTIFICATE-----`

const PrivateKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg7YwLqZ9Ny7o5fh50
5XvVe6pR9tTDBb3nGBRMWEcaQKqhRANCAARypEBgmw0HcPM5qsgRiqcoG7clO+cT
7MuqLMiFUwVmLGER+FUXZCsaKicYGWFuu1/E+b//NUbmrqKnmoyS0Iby
-----END PRIVATE KEY-----`

const PeerOrg1CA = `-----BEGIN CERTIFICATE-----
MIICJzCCAc6gAwIBAgIRAOt2caOGjH0/tS308yhB7tkwCgYIKoZIzj0EAwIwXjEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xDTALBgNVBAoTBG9yZzExEzARBgNVBAMTCnRsc2NhLm9yZzEwHhcN
MjAwNjIyMDAyNTAwWhcNMzAwNjIwMDAyNTAwWjBeMQswCQYDVQQGEwJVUzETMBEG
A1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzENMAsGA1UE
ChMEb3JnMTETMBEGA1UEAxMKdGxzY2Eub3JnMTBZMBMGByqGSM49AgEGCCqGSM49
AwEHA0IABCR8Vdcj3HS4h3b56lQ8PuNJLy+Q9+KmUbts/wm+bqruAJeF55IghlfZ
DLMNrdkCMzmcKNWN5NnYtCWEvZCNILSjbTBrMA4GA1UdDwEB/wQEAwIBpjAdBgNV
HSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zApBgNV
HQ4EIgQg794QgE2MOM6wPrkKEpY6LAJnCofYSJo+EZnzGTyl8pYwCgYIKoZIzj0E
AwIDRwAwRAIgeXGdlZ59qHg1k6wgjnfHxHHHsvukXFunccfZkkoK+DYCIDNr5e/+
12I9Exv0CgeSLrsd3AMGWNFLXsGeO9xePU7l
-----END CERTIFICATE-----`

const OrdererCA = `-----BEGIN CERTIFICATE-----
MIICQDCCAeagAwIBAgIRAO5u+03i/eKpuA1XD98uWpswCgYIKoZIzj0EAwIwajEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xEzARBgNVBAoTCmNvbnNvcnRpdW0xGTAXBgNVBAMTEHRsc2NhLmNv
bnNvcnRpdW0wHhcNMjAwNjIyMDAyNTAwWhcNMzAwNjIwMDAyNTAwWjBqMQswCQYD
VQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5j
aXNjbzETMBEGA1UEChMKY29uc29ydGl1bTEZMBcGA1UEAxMQdGxzY2EuY29uc29y
dGl1bTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABFW5ShgQaqPIGL/w58CK1SG5
eFXpAiX77u+ZTsbLbuaJ8Koid22mGeW55pNCbCcSiQwbPaaIZ8G6m9xBOKk2R4Oj
bTBrMA4GA1UdDwEB/wQEAwIBpjAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUH
AwEwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQg2+9b40BJCHhIXLmAK8XB/wjy
c/E6Dht75G66zyRDKyIwCgYIKoZIzj0EAwIDSAAwRQIhAPoenUeDyvwaUyyueb/b
GPQZVnUlqxaWjYU3u8r12zCbAiBxWf3c0ZBX5PvhZWXldNfnKdbt3aErH86BtAJd
0wHC5g==
-----END CERTIFICATE-----`

const channel = "trustchannel"

func main() {
	var clientConfig hlf.ClientConfig
	cryptoConfig := hlf.CryptoConfig{
		Family:    "ecdsa",
		Algorithm: "P256-SHA256",
		Hash:      "SHA2-256",
	}
	clientConfig.CryptoConfig = cryptoConfig

	clientConfig.EventPeers = map[string]hlf.PeerConfig{}
	clientConfig.Peers = map[string]hlf.PeerConfig{}
	clientConfig.Orderers = map[string]hlf.OrdererConfig{}

	clientConfig.EventPeers["peer0-org1"] = hlf.PeerConfig{"127.0.0.1:7051", true, PeerOrg1CA, "peer0"}
	clientConfig.Peers["peer0-org1"] = hlf.PeerConfig{"127.0.0.1:7051", true, PeerOrg1CA, "peer0"}
	clientConfig.Orderers["orderer1"] = hlf.OrdererConfig{"127.0.0.1:7050", true, OrdererCA, "orderer1"}

	admin, err := hlf.LoadCertFromBytes([]byte(Cert), []byte(PrivateKey))

	if err != nil {
		panic(fmt.Sprintf("Load Cert Error:%v", err))
	}
	admin.MspId = "Org1MSP"

	client, err := hlf.NewFabricClientFromConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	//block := queryBlockByNumber(client, admin)
	//
	//queryBlockByHash(client, admin, block.Header.Hash())

	//queryTxByHash(client, admin, "259df3b890543dc6c921ed38202e8a2ba41746241db300a0c89b4c856b159d21")

	//invoke(client, admin)
	query(client, admin)
}

func queryBlockByNumber(client *hlf.FabricClient, admin *hlf.Identity) *common.Block {
	blocks, err := client.QueryBlockByNumber(context.Background(), *admin, channel, 3, []string{"peer0-org1"})
	if err != nil {
		panic(err)
	}
	fmt.Println(blocks)
	return blocks[0]
}

func queryBlockByHash(client *hlf.FabricClient, admin *hlf.Identity, hash []byte) {
	blocks, err := client.QueryBlockByHash(context.Background(), *admin, channel, hash, []string{"peer0-org1"})
	if err != nil {
		panic(err)
	}

	fmt.Println(blocks)
}

func queryTxByHash(client *hlf.FabricClient, admin *hlf.Identity, hash string) {
	tx, err := client.QueryTransaction(context.Background(), *admin, channel, hash, []string{"peer0-org1"})
	if err != nil {
		panic(err)
	}
	fmt.Println(tx[0])
}

func invoke(client *hlf.FabricClient, admin *hlf.Identity) {
	gcc := hlf.ChainCode{
		ChannelId:    channel,
		Name:         "kvdb",
		Version:      "1.0",
		Type:         hlf.ChaincodeSpec_GOLANG,
		Args:         []string{"write", "k9", "100"},
		ArgBytes:     nil,
		TransientMap: nil,
	}
	resp, err := client.Invoke(context.Background(), *admin, gcc, []string{"peer0-org1"}, "orderer1")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func query(client *hlf.FabricClient, admin *hlf.Identity) {
	gcc := hlf.ChainCode{
		ChannelId:    channel,
		Name:         "kvdb",
		Version:      "1.0",
		Type:         hlf.ChaincodeSpec_GOLANG,
		Args:         []string{"read", "k9"},
		ArgBytes:     nil,
		TransientMap: nil,
	}
	resp, err := client.Query(context.Background(), *admin, gcc, []string{"peer0-org1"})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp[0])
}
