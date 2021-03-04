package db

import (
	"context"
	"strings"
	"time"

	"github.com/isnlan/coral/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func CreateUniqueIndex(ctx context.Context, coll *mongo.Collection, keys ...string) error {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := coll.Indexes()
	keysDoc := bsonx.Doc{}

	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// 创建索引
	_, err := indexView.CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		opts,
	)

	if err != nil {
		return errors.WithMessage(err, "EnsureIndex error")
	}
	return nil
}
