package hex

import (
	"github.com/ethereum/go-ethereum/common"
)

func Encode(b []byte) string {
	return common.Bytes2Hex(b)
}

func Decode(s string) []byte {
	return common.FromHex(s)
}
