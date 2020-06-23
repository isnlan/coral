package protos

//go:generate protoc -I protos/ protos/behavior.proto protos/auth.proto --go_out=plugins=grpc:protos
