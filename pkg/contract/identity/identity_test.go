package identity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentityToAddr(t *testing.T) {
	var cert = `-----BEGIN CERTIFICATE-----
MIIB/zCCAaWgAwIBAgIRAKaex32sim4PQR6kDPEPVnwwCgYIKoZIzj0EAwIwaTEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xFDASBgNVBAoTC2V4YW1wbGUuY29tMRcwFQYDVQQDEw5jYS5leGFt
cGxlLmNvbTAeFw0xNzA3MjYwNDM1MDJaFw0yNzA3MjQwNDM1MDJaMEoxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNp
c2NvMQ4wDAYDVQQDEwVwZWVyMDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABPzs
BSdIIB0GrKmKWn0N8mMfxWs2s1D6K+xvTvVJ3wUj3znNBxj+k2j2tpPuJUExt61s
KbpP3GF9/crEahpXXRajTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAA
MCsGA1UdIwQkMCKAIEvLfQX685pz+rh2q5yCA7e0a/a5IGDuJVHRWfp++HThMAoG
CCqGSM49BAMCA0gAMEUCIH5H9W3tsCrti6tsN9UfY1eeTKtExf/abXhfqfVeRChk
AiEA0GxTPOXVHo0gJpMbHc9B73TL5ZfDhujoDyjb8DToWPQ=
-----END CERTIFICATE-----`

	address, err := IntoAddress([]byte(cert))
	assert.NoError(t, err)

	assert.Equal(t, address.String(), "B3778BCEE2B9C349702E5832928730D2AED0AC07")

	getAddress, err := IntoIdentity([]byte(cert))
	assert.NoError(t, err)
	assert.Equal(t, string(getAddress), "b2d69605c52a7d05d07e9c5ccfe510dd02f9cfa41ba97144c5f6aeb38ca06d56")
}

func TestPublicKeyToAddr(t *testing.T) {
	pk := "9961335efcf2e2ab3f1a316b435719ff748220c8f83a1753fec429baf43600d6"
	addr, err := PublicKeyToAddr(pk)

	assert.NoError(t, err)
	fmt.Println(addr)
}
