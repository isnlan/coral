package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/snlansky/coral/pkg/net"
	"github.com/snlansky/coral/pkg/protos"
)

func main() {
	conn, err := net.New("192.168.2.2:8080")
	check(err)
	client, err := conn.Get()
	check(err)

	vmClient := protos.NewVMClient(client)
	ips, err := vmClient.GetNodeIps(context.Background(), &empty.Empty{})
	check(err)

	fmt.Println(ips.Ips)

	pods, err := vmClient.GetNamespacesPods(context.Background(), &protos.RequestNsPods{
		Ns:    "adminchain1org1",
		Label: "role",
		Filter: map[string]string{
			"ca":      "name",
			"orderer": "orderer-id",
			"peer":    "peer-id",
		},
	})
	check(err)
	fmt.Println(pods.Pods)

	status, err := vmClient.GetDeploymentStatus(context.Background(), &protos.RequestDeploymentStatus{
		Ns:   "kbcs",
		Name: "kchain-redis",
	})
	check(err)

	fmt.Println(status.Status)

	list, err := vmClient.GetDeploymentList(context.Background(), &protos.Namespaces{
		Ns: "kbcs",
	})
	check(err)

	fmt.Println(list.Deployments)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
