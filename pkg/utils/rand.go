package utils

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func MakeRandomString(len int) string {
	token := make([]byte, len)
	rand.Read(token)
	return hex.EncodeToString(token)
}
