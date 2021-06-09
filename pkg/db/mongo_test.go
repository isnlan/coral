package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/benweissmann/memongo"

	"github.com/stretchr/testify/assert"

	"github.com/smartystreets/goconvey/convey"
)

func TestMongoDaoImpl(t *testing.T) {
	mongoServer, err := memongo.Start("4.0.5")
	if err != nil {
		t.Fatal(err)
	}
	defer mongoServer.Stop()
	err = InitMongo(mongoServer.URI())
	assert.NoError(t, err)

	convey.Convey("test CreateUniqueIndex", t, func() {
		ctx := context.Background()

		mydb := GetDB("db1")

		coll := mydb.Collection("doc1")
		//err := CreateUniqueIndex(ctx, coll, "name")
		//convey.So(err, convey.ShouldBeNil)
		//
		//err1 := CreateUniqueIndex(ctx, coll, "name")
		//convey.So(err1, convey.ShouldBeNil)
		//
		//err2 := CreateUniqueIndex(ctx, coll, "age", "uid")
		//convey.So(err2, convey.ShouldBeNil)
		//
		//cursor, err := coll.Indexes().List(ctx)
		//convey.So(err, convey.ShouldBeNil)
		//defer cursor.Close(ctx)
		//
		////list := []interface{}{}
		//for cursor.Next(ctx) {
		//	var idx interface{}
		//	err := cursor.Decode(&idx)
		//	convey.So(err, convey.ShouldBeNil)
		//
		//	fmt.Println(idx)
		//}

		type User struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}
		var u User
		filter := bson.M(map[string]interface{}{})
		err := coll.FindOne(ctx, filter).Decode(&u)
		fmt.Println(errors.Is(err, mongo.ErrNoDocuments))

	})
}

func TestMongoApi(t *testing.T) {
	mongoServer, err := memongo.Start("4.0.5")
	if err != nil {
		t.Fatal(err)
	}
	defer mongoServer.Stop()

	client, err := MongoConnect("mongodb://admin:admin@127.0.0.1:27017/admin?authSource=admin")
	assert.NoError(t, err)

	coll := client.Database("blockchain_manager").Collection("api")

	ctx := context.Background()
	// id, err := primitive.ObjectIDFromHex("df05961e491bb6a77edeb7fc")
	// assert.NoError(t, err)

	var api0 API
	err = FindOne(ctx, coll, map[string]interface{}{"_id": primitive.NewObjectID()}, &api0)
	assert.True(t, errors.Is(err, mongo.ErrNoDocuments))

	// var api API
	// err = FindOne(ctx, coll, map[string]interface{}{"_id": id}, &api)
	// assert.NoError(t, err)
	// fmt.Println(api)

	var list []*API
	err = Find(ctx, coll, map[string]interface{}{}, 0, 0, primitive.D{}, func(cur *mongo.Cursor) error {
		var api API
		err := cur.Decode(&api)
		if err != nil {
			return err
		}

		list = append(list, &api)
		return nil
	})
	assert.NoError(t, err)

	fmt.Println(len(list))

	// apiNew := &API{
	// 	ID:             primitive.NewObjectID(),
	// 	Scheme:         "s",
	// 	Method:         "m",
	// 	Path:           "p",
	// 	AppName:        "",
	// 	APIName:        "",
	// 	APIType:        "",
	// 	DocURL:         "",
	// 	TotalCallCount: 0,
	// 	ClientCount:    0,
	// 	QPS:            0,
	// 	HistoryMaxQPS:  0,
	// 	SuccessCount:   0,
	// 	TotalLatency:   0,
	// }
	// err = InsertOne(ctx, coll, apiNew)
	// assert.NoError(t, err)

	// api.Scheme = "update"
	// err = UpdateOne(ctx, coll, map[string]interface{}{"_id": api.ID}, api)
	// assert.NoError(t, err)
	//
	// c, err := Count(ctx, coll, map[string]interface{}{})
	// assert.NoError(t, err)
	// fmt.Println(c)
	//
	// err = DeleteOne(ctx, coll, map[string]interface{}{"_id": api.ID})
	// assert.NoError(t, err)

	err = DeleteMany(ctx, coll, map[string]interface{}{"path": "/api/v2/orgs"})
	assert.NoError(t, err)

	var filter []map[string]interface{}
	filter = append(filter, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id": "$app_name",
			"count": map[string]interface{}{
				"$sum": 1,
			},
		},
	})

	m := map[string]interface{}{}
	res, err := Aggregations(ctx, coll, filter, reflect.TypeOf(m))
	assert.NoError(t, err)

	for _, item := range res {
		fmt.Println(item)
	}
}

type API struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`                            // 主键ID
	Scheme         string             `json:"scheme" bson:"scheme"`                     // http ws grpc
	Method         string             `json:"method" bson:"method"`                     // http method
	Path           string             `json:"path" bson:"path"`                         // http route path
	AppName        string             `json:"app_name" bson:"app_name"`                 // application name
	APIName        string             `json:"api_name" bson:"api_name"`                 // api 中文名称
	APIType        string             `json:"api_type" bson:"api_type"`                 // api 接口类型
	DocURL         string             `json:"doc_url" bson:"doc_url"`                   // 文档地址
	TotalCallCount int64              `json:"total_call_count" bson:"total_call_count"` // 累计调用数量
	ClientCount    int64              `json:"client_count" bson:"client_count"`         // 应用调用次数
	QPS            int64              `json:"qps" bson:"qps"`                           // 当前qps
	HistoryMaxQPS  int64              `json:"history_max_qps" bson:"history_max_qps"`   // 历史最高qps
	SuccessCount   int64              `json:"success_count" bson:"success_count"`       // 成功次数
	TotalLatency   int64              `json:"total_latency" bson:"total_latency"`       // 总耗时
}
