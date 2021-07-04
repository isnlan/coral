package main

import (
	"fmt"
	"os"

	"github.com/isnlan/coral/pkg/hlf"

	"context"
	"math/rand"
)

const ADM_PK = "/data/gopath/src/github.com/isnlan/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem"
const ADM_SK = "/data/gopath/src/github.com/isnlan/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/86454468ca20af00bd751106398310457f402b40b73f2e90bb07d46302227107_sk"

const channel = "mychannel"

const peer01 = "peer0.org1.example.com"
const peer11 = "peer1.org1.example.com"
const peer02 = "peer0.org2.example.com"
const peer12 = "peer1.org2.example.com"

const order0 = "orderer0"

const Org1MSP = "Org1MSP"
const Org2MSP = "Org2MSP"

func main() {
	// initialize Fabric client
	c, err := hlf.NewFabricClient("/data/gopath/src/github.com/isnlan/coral/test/hlf/client.yaml")
	if err != nil {
		fmt.Printf("Error loading file: %v", err)
		os.Exit(1)
	}

	// Fabric calls examples
	// some operations require admin certificate
	//createUpdateChannel(c)
	//joinChannel(c)
	//installCC(c)
	//instantiateCC(c)
	//queryInstalledChaincodes(c)
	//queryInstantiatedChaincodes(c)
	//queryChannels(c)
	//queryChannelInfo(c)
	//invoke(c, []string{"overwrite", "k1", "v1"})
	query(c)
	//queryTransaction(c, identity)
	//eventFullBlock(c, identity)
	//eventFilteredBlock(c, identity)
}

func eventFullBlock(client *hlf.FabricClient, identity *hlf.Identity) {
	ch := make(chan hlf.EventBlockResponse)
	ctx, cancel := context.WithCancel(context.Background())
	err := client.ListenForFullBlock(ctx, *identity, "peer0", "testchannel", ch)
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	for d := range ch {
		fmt.Println(d)
	}
}

func eventFilteredBlock(client *hlf.FabricClient, identity *hlf.Identity) {

	ch := make(chan hlf.EventBlockResponse)
	ctx, cancel := context.WithCancel(context.Background())
	err := client.ListenForFullBlock(ctx, *identity, "peer0", "testchannel", ch)
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	for d := range ch {
		fmt.Println(d)
	}
}

func invoke(client *hlf.FabricClient, q []string) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	admin.MspId = Org1MSP

	chaincode := hlf.ChainCode{
		ChannelId: channel,
		Type:      hlf.ChaincodeSpec_GOLANG,
		Name:      "kvdb",
		Version:   "1.0",
		Args:      q,
	}

	result, err := client.Invoke(context.Background(), *admin, chaincode, []string{peer01, peer11}, order0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(result)

}

func query(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	admin.MspId = Org1MSP

	chaincode := &hlf.ChainCode{
		ChannelId: channel,
		Type:      hlf.ChaincodeSpec_GOLANG,
		Name:      "kvdb",
		Version:   "1.0",
		Args:      []string{"read", "k1"},
	}

	result, err := client.Query(context.Background(), *admin, *chaincode, []string{peer01, peer11})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(result)
}

func queryTransaction(client *hlf.FabricClient, identity *hlf.Identity) {

	txid := "dd0945350a2e9e24515826f8fa6c7c8c5150001f0111478d7340d542dce6bd06"
	result, err := client.QueryTransaction(context.Background(), *identity, "testchannel", txid, []string{"peer0"})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(result)
}

func queryChannelInfo(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = Org1MSP
	result, err := client.QueryChannelInfo(context.Background(), *admin, channel, []string{peer01})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(result)
}

func queryChannels(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	admin.MspId = Org1MSP
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result, err := client.QueryChannels(context.Background(), *admin, []string{peer01})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(result)
}

func queryInstantiatedChaincodes(client *hlf.FabricClient) {

	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = Org1MSP

	result, err := client.QueryInstantiatedChainCodes(context.Background(), *admin, channel, []string{peer01, peer02})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(result)

}

func queryInstalledChaincodes(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = Org1MSP
	response, err := client.QueryInstalledChainCodes(context.Background(), *admin, []string{peer01, peer11, peer02})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(response)

}

func instantiateCC(client *hlf.FabricClient) {

	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = "Org1MSP"

	req := &hlf.ChainCode{
		Type:      hlf.ChaincodeSpec_GOLANG,
		ChannelId: channel,
		Name:      "kvdb",
		Version:   "1.1",
		Args:      []string{"init"}, // optional arguments for instantiation
	}

	// gohlf.CollectionConfig is new for v 1.1 and specify private collections for this chaincode. It is optional.

	cc := []hlf.CollectionConfig{}
	response, err := client.InstantiateChainCode(context.Background(), *admin, req, []string{peer01, peer11}, order0, "deploy", cc)
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}

	fmt.Println(response)
}

func installCC(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = Org1MSP
	req := &hlf.InstallRequest{
		ChainCodeType:    hlf.ChaincodeSpec_GOLANG,
		ChannelId:        channel,
		ChainCodeName:    "kvdb",
		ChainCodeVersion: "1.1",
		Namespace:        "kvdb",
		GoPath:           "/data/gopath",
	}
	response, err := client.InstallChainCode(context.Background(), *admin, req, []string{peer01, peer11})
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(response)
}

func joinChannel(client *hlf.FabricClient) {
	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = "Org1MSP"
	response, err := client.JoinChannel(context.Background(), *admin, channel, []string{"peer0.org1.example.com", "peer1.org1.example.com"}, "orderer0")
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	fmt.Println(response)

}

func createUpdateChannel(client *hlf.FabricClient) {

	admin, err := hlf.LoadCertFromFile(ADM_PK, ADM_SK)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Please note that we must provide MSPid manually because Identity is not from FabricCA
	admin.MspId = "Org1MSP"
	txpath := "/data/gopath/src/github.com/isnlan/e2e_cli/channel-artifacts/channel.tx"
	err = client.CreateUpdateChannel(*admin, txpath, "mychannel", "orderer0")
	fmt.Print(err)

}

func register(ca *hlf.FabricCAClient, identity *hlf.Identity) {

	// Optional list of attributes
	attr := []hlf.CaRegisterAttribute{{
		Name:  "option1",
		Value: "option1 value",
		ECert: true,
	},
		{
			Name:  "option2",
			Value: "option2 value",
			ECert: false,
		}}

	rr := hlf.CARegistrationRequest{
		EnrolmentId: "newUserName",
		Secret:      "qwerty",
		Affiliation: "comp1org",
		Type:        "user",
		Attrs:       attr}
	resp, err := ca.Register(identity, &rr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func reenroll(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	// optional attributes
	req := hlf.CaReEnrollmentRequest{
		Identity: identity,
		Attrs: []hlf.CaEnrollAttribute{
			{
				Name:     "option2",
				Optional: true,
			},
		},
	}
	resp, _, err := ca.ReEnroll(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func revoke(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	// To revoke user use:
	rr := hlf.CARevocationRequest{EnrollmentId: "newUser1"}

	// To revoke specific sertificate use:
	rr = hlf.CARevocationRequest{
		AKI:    "A84DEDAE57124E3D8305C9B8303E74A6EE196E27",
		Serial: "64e888fd586a6226016a70c22f2f5d95baa92599",
		GenCRL: true}
	r, err := ca.Revoke(identity, &rr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(r)

}

func getCaCerts(ca *hlf.FabricCAClient, identity *hlf.Identity) {

	resp, err := ca.GetCaCertificateChain("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)

}

func listAffiliation(ca *hlf.FabricCAClient, identity *hlf.Identity) {

	// path is optional
	resp, err := ca.ListAffiliations(identity, "organization1", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.CAName)
	fmt.Println(resp.Name)
	fmt.Println(resp.Affiliations)

}

func addAffiliation(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	req := hlf.CAAddAffiliationRequest{Name: "organization1.dep2", Force: false}
	resp, err := ca.AddAffiliation(identity, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func removeAffiliation(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	// CA must be configured to support affiliation removal
	req := hlf.CARemoveAffiliationRequest{Name: "organization1.department1", Force: false}
	resp, err := ca.RemoveAffiliation(identity, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func modifyAffiliation(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	req := hlf.CAModifyAffiliationRequest{Name: "organization1.department1", NewName: "org1.dep1", Force: true}
	resp, err := ca.ModifyAffiliation(identity, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func listAllIdentities(ca *hlf.FabricCAClient, identity *hlf.Identity) {

	resp, err := ca.ListAllIdentities(identity, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)

}

func removeIdentity(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	req := hlf.CARemoveIdentityRequest{Name: "newUser1", Force: false}
	resp, err := ca.RemoveIdentity(identity, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func modifyIdentity(ca *hlf.FabricCAClient, identity *hlf.Identity) {
	// see documentation for all fields that can be modified.
	req := hlf.CAModifyIdentityRequest{ID: "newUser1",
		Attributes: []hlf.CaRegisterAttribute{
			{
				Name:  "new1",
				ECert: true,
				Value: "new value 1",
			},
		},
		Secret: "new password"}
	resp, err := ca.ModifyIdentity(identity, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func getIdentity(ca *hlf.FabricCAClient, identity *hlf.Identity) {

	resp, err := ca.GetIdentity(identity, "newUser1", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp)

}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
