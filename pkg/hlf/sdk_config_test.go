package hlf

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSDKConfig(t *testing.T) {
	path := "/Users/snlan/go/src/github.com/isnlan/coral/pkg/hlf/testdata/connect.yaml"
	data, err := ioutil.ReadFile(path)
	assert.NoError(t, err)

	config, err := NewSDKConfig(data)
	assert.NoError(t, err)

	c, err := config.Into(true)
	assert.NoError(t, err)

	// sk := "/Users/snlan/nfs-data/baas/chains/adminchain1/artifacts/crypto-config/peerOrganizations/adminchain1org1/ca/priv_sk"
	// cert := "/Users/snlan/nfs-data/baas/chains/adminchain1/artifacts/crypto-config/peerOrganizations/adminchain1org1/ca/ca.adminchain1org1-cert.pem"

	cert1 := "/Users/snlan/nfs-data/baas/chains/adminchain1/artifacts/crypto-config/peerOrganizations/adminchain1org1/users/Admin@adminchain1org1/msp/signcerts/Admin@adminchain1org1-cert.pem"
	sk1 := "/Users/snlan/nfs-data/baas/chains/adminchain1/artifacts/crypto-config/peerOrganizations/adminchain1org1/users/Admin@adminchain1org1/msp/keystore/priv_sk"
	admin, err := LoadCertFromFile(cert1, sk1)
	assert.NoError(t, err)
	admin.MspId = "Org1MSP"

	client, err := NewFabricClientFromConfig(*c)
	assert.NoError(t, err)

	result, err := client.QueryChannels(context.Background(), *admin, []string{"peer0.adminchain1org1"})
	assert.NoError(t, err)

	fmt.Println(result)
	channel := "mychannel1"
	peer01 := "peer0.adminchain1org1"
	peer11 := "peer1.adminchain1org1"
	order0 := "orderer0.adminchain1orderer"

	if false {
		txpath := "/Users/snlan/nfs-data/baas/chains/adminchain1/artifacts/channel-artifacts/mychannel1.tx"
		err = client.CreateUpdateChannel(*admin, txpath, "mychannel1", order0)
		assert.NoError(t, err)
	}

	if false {
		resp, err := client.JoinChannel(context.Background(), *admin, channel, []string{peer01, peer11}, order0)
		assert.NoError(t, err)
		fmt.Println(resp)
	}

	if true {
		req := &InstallRequest{
			ChainCodeType:    ChaincodeSpec_GOLANG,
			ChannelId:        channel,
			ChainCodeName:    "kvdb",
			ChainCodeVersion: "1.1",
			Namespace:        "kvdb",
			GoPath:           "/Users/snlan/go",
		}
		response, err := client.InstallChainCode(context.Background(), *admin, req, []string{peer01, peer11})
		assert.NoError(t, err)
		fmt.Println(response)
	}

	if false {
		req := &ChainCode{
			Type:      ChaincodeSpec_GOLANG,
			ChannelId: channel,
			Name:      "kvdb",
			Version:   "1.1",
			Args:      []string{"init"}, // optional arguments for instantiation
		}

		// gohlf.CollectionConfig is new for v 1.1 and specify private collections for this chaincode. It is optional.

		cc := []CollectionConfig{}
		response, err := client.InstantiateChainCode(context.Background(), *admin, req, []string{peer01, peer11}, order0, "deploy", cc)
		assert.NoError(t, err)
		fmt.Println(response)
	}
}
