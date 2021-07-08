package hlf

import (
	"encoding/pem"
	"fmt"
	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/hlf/crypto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadECCertFromBytes(t *testing.T) {
	const Cert = `-----BEGIN CERTIFICATE-----
MIIChzCCAi6gAwIBAgIUWF+TnGBBFMTBEX4tHscrwB26sNkwCgYIKoZIzj0EAwIw
WzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xDTALBgNVBAoTBG9yZzExEDAOBgNVBAMTB2NhLm9yZzEwHhcN
MjAwNjIzMDI1OTAwWhcNMjEwNjIzMDMwNDAwWjAxMRowCwYDVQQLEwR1c2VyMAsG
A1UECxMEb3JnMTETMBEGA1UEAxMKMTU5Mjg4MTQ1NDBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABH8X5EkL0+pkIByAUiHUmAsuE9sTI/yXPfCj5o4Pxer4EC9ihYKR
VqYYJrz/EqsdDmr908mfdgC5SBj3899LzUGjgfkwgfYwDgYDVR0PAQH/BAQDAgeA
MAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFIkfsJNLiJLDbpqVDEo0Jv8YYHYqMCsG
A1UdIwQkMCKAIPi3Pb5wquira1i886vqLO17BxyEJZASbouE9qy+7uFsMBkGA1Ud
EQQSMBCCDjEyNy4wLjAuMTo3MDU0MG8GCCoDBAUGBwgBBGN7ImF0dHJzIjp7IjE1
OTI4ODE0NTQiOiIiLCJoZi5BZmZpbGlhdGlvbiI6Im9yZzEiLCJoZi5FbnJvbGxt
ZW50SUQiOiIxNTkyODgxNDU0IiwiaGYuVHlwZSI6InVzZXIifX0wCgYIKoZIzj0E
AwIDRwAwRAIgX6i/fq/UE5+mxzZMRtcyA0HRD8/yxeHDquCk7ixGVscCIBkMHQHt
fYZ09AOH/MZGjMjI8Ndyd111CbFSxuxAKFDR
-----END CERTIFICATE-----`

	const PrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKrKuuc5n1iVPRDErruqft9D/a0GSXz+PLrO1vsXn3iqoAoGCCqGSM49
AwEHoUQDQgAEfxfkSQvT6mQgHIBSIdSYCy4T2xMj/Jc98KPmjg/F6vgQL2KFgpFW
phgmvP8Sqx0Oav3TyZ92ALlIGPfz30vNQQ==
-----END EC PRIVATE KEY-----`
	identity, err := LoadECCertFromBytes([]byte(Cert), []byte(PrivateKey))
	assert.NoError(t, err)
	addr, err := identity.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, addr.String(), "97362A7B592DAF95AFD0B046551EC9E62BE7607F")

	idBytes :=  pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: identity.Certificate.Raw})

	expirationTime := crypto.ExpiresAt(idBytes)
	if !expirationTime.IsZero() && time.Now().After(expirationTime) {
		err := errors.New("proposal client identity expired")
		assert.NoError(t, err)
	}
	fmt.Println(expirationTime.IsZero())
}
