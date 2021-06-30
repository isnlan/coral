package offchain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	qs := New("http://127.0.0.1:8089")
	info, err := qs.QueryChannelInfo(context.Background(), "60dc61690c779df8628c051e", "channel")

	assert.NoError(t, err)
	fmt.Println(info)

	{
		var filter []map[string]interface{}
		match := map[string]interface{}{
			"number": 1,
		}

		filter = append(filter, map[string]interface{}{
			"$match": match,
		})

		blocks, err := qs.QueryBlocks(context.Background(), "60dc61690c779df8628c051e", "channel", filter)
		assert.NoError(t, err)
		fmt.Println(blocks[0])

	}

	{
		var filter []map[string]interface{}
		match := map[string]interface{}{
			"tx_id": "ee18f6cab8b6395e04220135f191ae6e5164a9c56b0786c8d89ca44848630652",
		}

		filter = append(filter, map[string]interface{}{
			"$match": match,
		})

		txs, err := qs.QueryTxs(context.Background(), "60dc61690c779df8628c051e", "channel", filter)
		assert.NoError(t, err)
		fmt.Println(txs[0])

	}

}
