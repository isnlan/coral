package protos

//go:generate protoc --go_out=plugins=grpc:. model.proto behavior.proto auth.proto vm.proto network.proto
