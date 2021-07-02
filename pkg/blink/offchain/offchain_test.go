package offchain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	qs := New("http://127.0.0.1:8089")
	chain := "60dc61690c779df8628c051e"
	channel := "channel"
	info, err := qs.QueryChannelInfo(context.Background(), chain, channel)

	assert.NoError(t, err)
	fmt.Println(info)

	{
		match := map[string]interface{}{
			"number": 1,
		}

		blocks, count, err := qs.QueryBlocks(context.Background(), chain, channel, match, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, count, int64(1))
		fmt.Println(blocks[0])

	}

	{
		match := map[string]interface{}{
			"channel_id": channel,
		}

		txs, count, err := qs.QueryTxs(context.Background(), chain, channel, match, 0, 1)
		assert.NoError(t, err)
		// assert.Equal(t, count, int64(1))
		fmt.Println(count)
		fmt.Println(len(txs))
	}
}
