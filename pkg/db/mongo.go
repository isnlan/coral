package db

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _mongoClient *mongo.Client

func InitMongo(uri string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return errors.WithMessage(err, "mongo connect error")
	}
	err = _mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.WithMessage(err, "mongo ping error")
	}

	return nil
}

func GetMongoClient() *mongo.Client {
	return _mongoClient
}

func GetDB(name string) *mongo.Database {
	return _mongoClient.Database(name)
}

type Dao struct {
	coll *mongo.Collection
}

func NewDaoByColl(coll *mongo.Collection) *Dao {
	return &Dao{coll: coll}
}

func NewDao(db *mongo.Database, table string) *Dao {
	return &Dao{coll: db.Collection(table)}
}

func (d *Dao) Save(ctx context.Context, doc interface{}) error {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	_, err := d.coll.InsertOne(ctx, doc)
	return err
}

func (d *Dao) FindOne(ctx context.Context, query map[string]interface{}, doc interface{}) (bool, error) {
	filter := bson.M(query)
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	err := d.coll.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *Dao) UpdateOne(ctx context.Context, condition map[string]interface{}, doc interface{}) error {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	filter := bson.M(condition)
	update := bson.D{
		{Key: "$set", Value: doc},
	}
	_, err := d.coll.UpdateOne(ctx, filter, update)
	return err
}

func (d *Dao) DeleteOne(ctx context.Context, condition map[string]interface{}) error {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	filter := bson.M(condition)
	_, err := d.coll.DeleteOne(ctx, filter)
	return err
}

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
