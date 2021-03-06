package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func Generate() snowflake.ID {
	// Generate a snowflake ID.
	return node.Generate()
}
