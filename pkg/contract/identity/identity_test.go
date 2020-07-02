package identity

import (
	"encoding/json"
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

func TestAddress_Bytes(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIICjDCCAjKgAwIBAgIUTqNMzlIsNdUS+ElmBDsNZjkk+EAwCgYIKoZIzj0EAwIw
WzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xDTALBgNVBAoTBG9yZzExEDAOBgNVBAMTB2NhLm9yZzEwHhcN
MjAwNjMwMDgzNTAwWhcNMjEwNjMwMDg0MDAwWjAzMRwwDQYDVQQLEwZjbGllbnQw
CwYDVQQLEwRvcmcxMRMwEQYDVQQDEwoxNTkzNTA2NDExMFkwEwYHKoZIzj0CAQYI
KoZIzj0DAQcDQgAE+/7qRd3ogoUs0pSYItN5tjxWQCRBIou2mhZm8I7Rj0YyaB0e
2S5uOCrjxMVQgSHXKfwWfxdkf4S8ck8NJfqn1KOB+zCB+DAOBgNVHQ8BAf8EBAMC
B4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUsEim9OGf1L/78XkaDQAYjZHglw8w
KwYDVR0jBCQwIoAg+Lc9vnCq6KtrWLzzq+os7XsHHIQlkBJui4T2rL7u4WwwGQYD
VR0RBBIwEIIOMTI3LjAuMC4xOjcwNTQwcQYIKgMEBQYHCAEEZXsiYXR0cnMiOnsi
MTU5MzUwNjQxMSI6IiIsImhmLkFmZmlsaWF0aW9uIjoib3JnMSIsImhmLkVucm9s
bG1lbnRJRCI6IjE1OTM1MDY0MTEiLCJoZi5UeXBlIjoiY2xpZW50In19MAoGCCqG
SM49BAMCA0gAMEUCIQDYr6L2eqJNBgCE1a63Bwro5EcniyQG7pg1RBNYGQqs2wIg
MXXw8KMx+fR4BbSs2ra7Jiga3PHYXs/B9HHy0TDtD6o=
-----END CERTIFICATE-----`
	address, err := IntoAddress([]byte(cert))
	assert.NoError(t, err)
	fmt.Println(address.String())
	str := address.String()
	naddr := MustAddressFromHexString(str)
	assert.Equal(t, address, naddr)

	n2addr, err := AddressFromHexString(str)
	assert.NoError(t, err)
	assert.Equal(t, address, n2addr)

	n3addr, err := AddressFromBytes(address.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, address, n3addr)
}

func TestAddress_MarshalJSON(t *testing.T) {
	type User struct {
		Name   string             `json:"name"`
		Addr   Address            `json:"addr"`
		Family map[Address]string `json:"family"`
	}
	type Human struct {
		Name   string            `json:"name"`
		Addr   string            `json:"addr"`
		Family map[string]string `json:"family"`
	}

	cert := `-----BEGIN CERTIFICATE-----
MIICjDCCAjKgAwIBAgIUTqNMzlIsNdUS+ElmBDsNZjkk+EAwCgYIKoZIzj0EAwIw
WzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xDTALBgNVBAoTBG9yZzExEDAOBgNVBAMTB2NhLm9yZzEwHhcN
MjAwNjMwMDgzNTAwWhcNMjEwNjMwMDg0MDAwWjAzMRwwDQYDVQQLEwZjbGllbnQw
CwYDVQQLEwRvcmcxMRMwEQYDVQQDEwoxNTkzNTA2NDExMFkwEwYHKoZIzj0CAQYI
KoZIzj0DAQcDQgAE+/7qRd3ogoUs0pSYItN5tjxWQCRBIou2mhZm8I7Rj0YyaB0e
2S5uOCrjxMVQgSHXKfwWfxdkf4S8ck8NJfqn1KOB+zCB+DAOBgNVHQ8BAf8EBAMC
B4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUsEim9OGf1L/78XkaDQAYjZHglw8w
KwYDVR0jBCQwIoAg+Lc9vnCq6KtrWLzzq+os7XsHHIQlkBJui4T2rL7u4WwwGQYD
VR0RBBIwEIIOMTI3LjAuMC4xOjcwNTQwcQYIKgMEBQYHCAEEZXsiYXR0cnMiOnsi
MTU5MzUwNjQxMSI6IiIsImhmLkFmZmlsaWF0aW9uIjoib3JnMSIsImhmLkVucm9s
bG1lbnRJRCI6IjE1OTM1MDY0MTEiLCJoZi5UeXBlIjoiY2xpZW50In19MAoGCCqG
SM49BAMCA0gAMEUCIQDYr6L2eqJNBgCE1a63Bwro5EcniyQG7pg1RBNYGQqs2wIg
MXXw8KMx+fR4BbSs2ra7Jiga3PHYXs/B9HHy0TDtD6o=
-----END CERTIFICATE-----`
	address, err := IntoAddress([]byte(cert))
	assert.NoError(t, err)
	u := User{Name: "snlan", Addr: address, Family: map[Address]string{address: "my"}}
	bytes, err := json.Marshal(&u)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
	var u2 User
	err = json.Unmarshal(bytes, &u2)
	assert.NoError(t, err)
	fmt.Println(u2)
	assert.Equal(t, u, u2)

	var h Human
	err = json.Unmarshal(bytes, &h)
	assert.NoError(t, err)
	fmt.Println(h)
}

func TestAddress_MarshalText(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIICjDCCAjKgAwIBAgIUTqNMzlIsNdUS+ElmBDsNZjkk+EAwCgYIKoZIzj0EAwIw
WzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xDTALBgNVBAoTBG9yZzExEDAOBgNVBAMTB2NhLm9yZzEwHhcN
MjAwNjMwMDgzNTAwWhcNMjEwNjMwMDg0MDAwWjAzMRwwDQYDVQQLEwZjbGllbnQw
CwYDVQQLEwRvcmcxMRMwEQYDVQQDEwoxNTkzNTA2NDExMFkwEwYHKoZIzj0CAQYI
KoZIzj0DAQcDQgAE+/7qRd3ogoUs0pSYItN5tjxWQCRBIou2mhZm8I7Rj0YyaB0e
2S5uOCrjxMVQgSHXKfwWfxdkf4S8ck8NJfqn1KOB+zCB+DAOBgNVHQ8BAf8EBAMC
B4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUsEim9OGf1L/78XkaDQAYjZHglw8w
KwYDVR0jBCQwIoAg+Lc9vnCq6KtrWLzzq+os7XsHHIQlkBJui4T2rL7u4WwwGQYD
VR0RBBIwEIIOMTI3LjAuMC4xOjcwNTQwcQYIKgMEBQYHCAEEZXsiYXR0cnMiOnsi
MTU5MzUwNjQxMSI6IiIsImhmLkFmZmlsaWF0aW9uIjoib3JnMSIsImhmLkVucm9s
bG1lbnRJRCI6IjE1OTM1MDY0MTEiLCJoZi5UeXBlIjoiY2xpZW50In19MAoGCCqG
SM49BAMCA0gAMEUCIQDYr6L2eqJNBgCE1a63Bwro5EcniyQG7pg1RBNYGQqs2wIg
MXXw8KMx+fR4BbSs2ra7Jiga3PHYXs/B9HHy0TDtD6o=
-----END CERTIFICATE-----`
	address, err := IntoAddress([]byte(cert))
	assert.NoError(t, err)
	text, err := address.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, string(text), address.String())

	var n1addr Address
	err = n1addr.UnmarshalText([]byte(address.String()))
	assert.NoError(t, err)
	assert.Equal(t, n1addr, address)
}

func TestNewRandomAddress(t *testing.T) {
	address := NewRandomAddress()
	fmt.Println(address)
}

func TestAddress_UnmarshalJSON(t *testing.T) {
	type User struct {
		Name string   `json:"name"`
		Addr *Address `json:"addr,omitempty"` //`json:"addr,omitempty"`
	}
	type Human struct {
		Name string `json:"name"`
		Addr string `json:"addr"`
	}

	str := `{"name":"snlan","addr":""}`

	var u User
	err := json.Unmarshal([]byte(str), &u)
	assert.NoError(t, err)
	fmt.Println("user", u)
	fmt.Println("user nil", User{Name: "s1"})

	fmt.Println("-----")
	marshal, err := json.Marshal(&User{Name: "lucy", Addr: nil})
	assert.NoError(t, err)
	fmt.Println(string(marshal))
}
