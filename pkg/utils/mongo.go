package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MakeMongoIdFromString(str string) string {
	bytes := sha256.Sum256([]byte(str))
	var dst [12]byte
	copy(dst[:], bytes[:])
	return hex.EncodeToString(dst[:])
}

func MakeMongoIdf(format string, a ...interface{}) string {
	return MakeMongoIdFromString(fmt.Sprintf(format, a...))
}
