package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	err := InitMongo("mongodb://root:uMKxr0EClp@127.0.0.1:27017/admin?authSource=admin")
	if err != nil {
		panic(err)
	}
}
func TestDao_UpdateOne(t *testing.T) {
	dao := NewDao(GetDB("mytest"), "testcoll")

	err := dao.Save(context.Background(), map[string]interface{}{"id": 1, "value": 13})
	assert.NoError(t, err)

	filter := map[string]interface{}{"id": 1}
	doc := map[string]interface{}{"id": 1, "value": 15}
	err = dao.UpdateOne(context.Background(), filter, doc)
	assert.NoError(t, err)

	err = dao.DeleteOne(context.Background(), map[string]interface{}{"id": 2})
	assert.NoError(t, err)
}
