package contract

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/isnlan/coral/pkg/contract"
	"github.com/isnlan/coral/pkg/contract/impl"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserService struct {
}

func (s *UserService) Add(stub contract.IContractStub, user *User) int {
	fmt.Println(user)
	b, err := json.Marshal(user)
	contract.Check(err, contract.ERR_JSON_MARSHAL)
	err = stub.PutState(user.Id, b)
	contract.Check(err, contract.ERR_INTERNAL_INVALID)
	return 0
}

func (s *UserService) Get(stub contract.IContractStub, id string) *User {
	state, err := stub.GetState(id)
	contract.Check(err, contract.ERR_INTERNAL_INVALID)
	var u User
	err = json.Unmarshal(state, &u)
	contract.Check(err)

	return &u
}

func (s *UserService) TryPanic(stub contract.IContractStub, trigger bool) int {
	if trigger {
		contract.Throw(contract.ERR_RUNTIME)
	}
	return 0
}

func TestNewFabricChaincode(t *testing.T) {
	cc := impl.NewFabricChaincode()
	cc.Register(&UserService{})

	stub := NewMockStub("UserCC", cc)

	stub.MockTransactionStart("init")

	var parems []interface{}
	var resp pb.Response

	resp = stub.MockInit("1", [][]byte{[]byte("init")})
	assert.Equal(t, resp.Payload, []byte("SUCCESS"))

	parems = []interface{}{&User{
		Id:   "user1",
		Name: "snlan",
	}}
	buf, err := json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("UserService.Add"), buf})
	v, _ := json.Marshal(0)
	assert.Equal(t, resp.Payload, v)

	parems = []interface{}{
		"user1",
	}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("3", [][]byte{[]byte("UserService.Get"), buf})
	fmt.Println(string(resp.Payload))

	parems = []interface{}{
		true,
	}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("4", [][]byte{[]byte("UserService.TryPanic"), buf})
	fmt.Println(resp.Message)
	fmt.Println(string(resp.Payload))
}
