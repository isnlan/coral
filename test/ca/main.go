package main

import (
	"fmt"
	"time"

	"github.com/isnlan/coral/pkg/hlf"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	conf := hlf.CAConfig{
		CryptoConfig: hlf.CryptoConfig{
			Family:    "ecdsa",
			Algorithm: "P256-SHA256",
			Hash:      "SHA2-256",
		},
		Uri:               "http://127.0.0.1:7054",
		SkipTLSValidation: true,
		MspId:             "Org1MSP",
	}
	ca, err := hlf.NewCaClientFromConfig(conf, nil)
	check(err)

	admin, cert, err := ca.Enroll(hlf.CaEnrollmentRequest{
		EnrollmentId: "admin",
		Secret:       "adminpw",
		Profile:      "",
		Label:        "",
		CAName:       "",
		Hosts:        nil,
		Attrs:        nil,
	})
	check(err)
	fmt.Println(admin)
	fmt.Println(string(cert))
	address, err := admin.GetAddress()
	check(err)
	fmt.Println(address.String())
	fmt.Println("--------1-------")
	{
		name := fmt.Sprintf("snlan%d", time.Now().Unix())
		attr := []hlf.CaRegisterAttribute{{
			Name:  name,
			Value: "value1",
			ECert: true,
		}}
		rr := &hlf.CARegistrationRequest{
			EnrolmentId:    name,
			Type:           "client", // 这里只能"must be a client, a peer, an orderer or an admin",因为crypto-config.yaml 设置了EnableNodeOUs: true
			Secret:         "mypasswd",
			MaxEnrollments: 0,
			Affiliation:    "org1",
			Attrs:          attr,
			CAName:         "",
		}

		secret, err := ca.Register(admin, rr)
		fmt.Println("--------2-------")
		check(err)
		fmt.Println("--", secret)

		user, cert, err := ca.Enroll(hlf.CaEnrollmentRequest{
			EnrollmentId: name,
			Secret:       "mypasswd",
			Profile:      "",
			Label:        "",
			CAName:       "",
			Hosts:        nil,
			Attrs:        nil,
		})

		fmt.Println("--------3-------")

		check(err)
		fmt.Println(user)
		fmt.Println(string(cert))
	}
}
