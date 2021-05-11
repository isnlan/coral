package ca

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/isnlan/coral/pkg/hlf"
	"github.com/stretchr/testify/assert"
)

func TestCA_Enroll(t *testing.T) {
	CryptoConfig := hlf.CryptoConfig{
		Family:    "ecdsa",
		Algorithm: "P256-SHA256",
		Hash:      "SHA2-256",
	}
	crypto, err := hlf.NewECCryptSuiteFromConfig(CryptoConfig)
	assert.NoError(t, err)

	key, err := crypto.GenerateKey()
	assert.NoError(t, err)

	conf := &Config{
		ParentServerURL: "",
		CertFile:        "/Users/snlan/go/src/github.com/isnlan/coral/test/ca/ca.org1-cert.pem",
		KeyFile:         "/Users/snlan/go/src/github.com/isnlan/coral/test/ca/key.pem",
	}

	csr, err := crypto.CreateCertificateRequest("myname", key, []string{"127.0.0.1:7054"})
	assert.NoError(t, err)
	ca, err := New(conf)
	assert.NoError(t, err)
	req := MakeEmptySignRequest(csr)
	cert, err := ca.Enroll(req, "myname", "mupass")
	assert.NoError(t, err)
	fmt.Println(string(csr), string(cert))

	a, _ := pem.Decode(cert)
	certObj, err := x509.ParseCertificate(a.Bytes)
	assert.NoError(t, err)
	id := &hlf.Identity{Certificate: certObj, PrivateKey: key, MspId: "Org1MSP"}

	address, err := id.GetAddress()
	assert.NoError(t, err)
	fmt.Println(address)
}
