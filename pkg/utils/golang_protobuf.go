package utils

import "github.com/golang/protobuf/proto"

func GolangProtobufEqual(x, y proto.Message) bool {
	return proto.Equal(x, y)
}
