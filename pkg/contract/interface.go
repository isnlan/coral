package contract

import (
	"time"

	"github.com/snlansky/coral/pkg/contract/identity"
)

type IContractStub interface {
	GetArgs() [][]byte
	GetTxID() string
	GetChannelID() string
	GetAddress() (identity.Address, error)
	GetState(key string) ([]byte, error)
	PutState(key string, value []byte) error
	DelState(key string) ([]byte, error)
	CreateCompositeKey(objectType string, attributes []string) (string, error)
	SplitCompositeKey(compositeKey string) (string, []string, error)
	GetTxTimestamp() (time.Time, error)
	SetEvent(name string, payload []byte) error
	InvokeContract(contractName string, args [][]byte, channel string) ([]byte, error)
	GetOriginStub() interface{}
}
