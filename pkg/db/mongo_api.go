package db

import (
	"context"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/isnlan/coral/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsDup(err error) bool {
	merr, ok := err.(mongo.WriteException)
	if ok {
		for _, we := range merr.WriteErrors {
			if we.Code == 11000 ||
				we.Code == 11001 ||
				we.Code == 12582 ||
				(we.Code == 16460 && strings.Contains(we.Message, " E11000 ")) {
				return true
			}
		}
	}

	return false
}

func IsConnectionError(err error) bool {
	rawErr := errors.Cause(err)
	if rawErr == nil {
		return false
	}

	msg := rawErr.Error()

	switch {
	case strings.Contains(msg, "connection is closed"):
		return true
	case strings.Contains(msg, "failed to read"):
		return true
	case strings.Contains(msg, "failed to set read deadline"):
		return true
	case strings.Contains(msg, "incomplete read of message header"):
		return true
	case strings.Contains(msg, "incomplete read of full message"):
		return true
	default:
		return false
	}
}

func InsertOne(ctx context.Context, coll *mongo.Collection, data interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	_, err := coll.InsertOne(ctx, data)

	return err
}

type Decoder interface {
	Decode(val interface{}) error
}

func Find(ctx context.Context, coll *mongo.Collection, condition map[string]interface{}, limit, skip int64, sort bson.D, f func(decoder Decoder) error) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	filter := bson.M(condition)

	cur, err := coll.Find(ctx, filter, &options.FindOptions{Limit: &limit, Skip: &skip, Sort: sort})
	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := f(cur)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindOne :must check mongo.ErrNoDocuments
func FindOne(ctx context.Context, coll *mongo.Collection, query map[string]interface{}, v interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	return coll.FindOne(ctx, bson.M(query)).Decode(v)
}

func Count(ctx context.Context, coll *mongo.Collection, condition map[string]interface{}) (int64, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	count, err := coll.CountDocuments(ctx, bson.M(condition))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateOne(ctx context.Context, coll *mongo.Collection, condition map[string]interface{}, v interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	update := bson.D{
		bson.E{Key: "$set", Value: v},
	}

	_, err := coll.UpdateOne(ctx, bson.M(condition), update)

	return err
}

func DeleteOne(ctx context.Context, coll *mongo.Collection, condition map[string]interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	_, err := coll.DeleteOne(ctx, bson.M(condition))

	return err
}

func DeleteMany(ctx context.Context, coll *mongo.Collection, condition map[string]interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	_, err := coll.DeleteMany(ctx, bson.M(condition))

	return err
}

func Aggregations(ctx context.Context, coll *mongo.Collection, pipeline interface{}, t reflect.Type) ([]interface{}, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer cursor.Close(ctx)

	list := []interface{}{}

	for cursor.Next(ctx) {
		d := reflect.New(t)
		err := cursor.Decode(d.Interface())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		list = append(list, d.Interface())
	}

	return list, nil
}
