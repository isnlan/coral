package ca

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/isnlan/coral/pkg/hlf"
	crypto2 "github.com/isnlan/coral/pkg/hlf/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCA_Enroll(t *testing.T) {
	CryptoConfig := hlf.CryptoConfig{
		Family:    "ecdsa",
		Algorithm: "P256-SHA256",
		Hash:      "SHA2-256",
	}

	conf := &Config{
		ParentServerURL: "",
		CertFile:        "/data/gopath/src/github.com/isnlan/coral/test/ca/ca.org1-cert.pem",
		KeyFile:         "/data/gopath/src/github.com/isnlan/coral/test/ca/key.pem",
	}

	crypto, err := hlf.NewECCryptSuiteFromConfig(CryptoConfig)
	assert.NoError(t, err)

	key, err := crypto.GenerateKey()
	assert.NoError(t, err)

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
	fmt.Println(1, address)


	idBytes :=  pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: id.Certificate.Raw})


	expirationTime := crypto2.ExpiresAt(idBytes)
	if !expirationTime.IsZero() && time.Now().After(expirationTime) {
		fmt.Println("proposal client identity expired")
	}
	fmt.Println(expirationTime.IsZero())
}