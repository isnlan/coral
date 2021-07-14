package metrics

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScraper(t *testing.T) {
	list, err := ScrapeMetrics(context.Background(), "http://127.0.0.1:9001/metrics")
	assert.NoError(t, err)
	for _, v := range list {
		fmt.Printf("%+#v\n", v)
	}
}
