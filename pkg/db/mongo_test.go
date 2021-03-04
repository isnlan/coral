package db

import (
	"context"
	"fmt"
	"testing"

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
		err := CreateUniqueIndex(ctx, coll, "name")
		convey.So(err, convey.ShouldBeNil)

		err1 := CreateUniqueIndex(ctx, coll, "name")
		convey.So(err1, convey.ShouldBeNil)

		err2 := CreateUniqueIndex(ctx, coll, "age", "uid")
		convey.So(err2, convey.ShouldBeNil)

		cursor, err := coll.Indexes().List(ctx)
		convey.So(err, convey.ShouldBeNil)
		defer cursor.Close(ctx)

		//list := []interface{}{}
		for cursor.Next(ctx) {
			var idx interface{}
			err := cursor.Decode(&idx)
			convey.So(err, convey.ShouldBeNil)

			fmt.Println(idx)
		}

	})

}
