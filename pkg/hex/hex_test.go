package hex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	b2 := Decode("0xec9272bcbc7ffc07ef35b6be2f0e26a2aeb11dd6cb55251bd38a8ddc1e7913b0")
	b3 := Decode("ec9272bcbc7ffc07ef35b6be2f0e26a2aeb11dd6cb55251bd38a8ddc1e7913b0")

	assert.Equal(t, b2, b3)
}
