package identity

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/tmthrgd/go-hex"
	"golang.org/x/crypto/sha3"
)

// Deprecated. should use IntoAddress
func IntoIdentity(creatorByte []byte) ([]byte, error) {
	certStart := bytes.Index(creatorByte, []byte("-----BEGIN"))
	if certStart == -1 {
		return nil, errors.New("no creator certificate found")
	}
	certText := creatorByte[certStart:]

	bl, _ := pem.Decode(certText)
	if bl == nil {
		return nil, errors.New("could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return nil, errors.New("parse Certificate failed")
	}

	if pub, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
		pubKey := append(pub.X.Bytes(), pub.Y.Bytes()...)
		publicSHA256 := sha256.Sum256(pubKey)
		address := hex.EncodeToString(publicSHA256[:])

		return []byte(address), nil
	}

	return nil, errors.New("only support ECDSA")
}

// IntoAddress computes 160 bits address from the public key encoded in an identity.
func IntoAddress(creatorByte []byte) (Address, error) {
	certStart := bytes.Index(creatorByte, []byte("-----BEGIN"))
	if certStart == -1 {
		return ZeroAddress, errors.New("no creator certificate found")
	}
	certText := creatorByte[certStart:]

	bl, _ := pem.Decode(certText)
	if bl == nil {
		return ZeroAddress, errors.New("could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return ZeroAddress, fmt.Errorf("failed to parse certificate: %s", err)
	}

	return PublicKeyIntoAddress(cert.PublicKey)
}

func PublicKeyIntoAddress(pk interface{}) (Address, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return ZeroAddress, fmt.Errorf("unable to marshal public key: %s", err)
	}

	// We want the last 160 bits of the sha3-256 sum as the address
	sum := sha3.Sum256(pubkeyBytes)
	return AddressFromBytes(sum[12:])
}

// ed25519 public key to address
func PublicKeyToAddr(pk string) (Address, error) {
	pubkeyBytes, err := hex.DecodeString(pk)
	if err != nil {
		return ZeroAddress, fmt.Errorf("unable to marshal public key: %s", err)
	}

	// We want the last 160 bits of the sha3-256 sum as the address
	sum := sha3.Sum256(pubkeyBytes)
	return AddressFromBytes(sum[12:])
}
