package impl

import (
	"errors"
	"fmt"
	"time"

	"github.com/snlansky/coral/pkg/contract"
	"github.com/snlansky/coral/pkg/contract/identity"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type FabricContractStub struct {
	stub    shim.ChaincodeStubInterface
	creator func() []byte
}

func NewFabricContractStub(stub shim.ChaincodeStubInterface) contract.IContractStub {
	return &FabricContractStub{stub: stub}
}

func (f *FabricContractStub) setCreatorFactory(creator func() []byte) {
	f.creator = creator
}

func (f *FabricContractStub) GetArgs() [][]byte {
	return f.stub.GetArgs()
}

func (f *FabricContractStub) GetTxID() string {
	return f.stub.GetTxID()
}

func (f *FabricContractStub) GetChannelID() string {
	return f.stub.GetChannelID()
}

func (f *FabricContractStub) GetAddress() ([]byte, error) {
	creatorByte, err := f.stub.GetCreator()
	if err != nil {
		return nil, err
	}
	return identity.IntoIdentity(creatorByte)
}

func (f *FabricContractStub) GetState(key string) ([]byte, error) {
	return f.stub.GetState(key)
}

func (f *FabricContractStub) PutState(key string, value []byte) error {
	return f.stub.PutState(key, value)
}

func (f *FabricContractStub) DelState(key string) ([]byte, error) {
	buf, err := f.stub.GetState(key)
	if err != nil {
		return nil, err
	}
	err = f.stub.DelState(key)
	return buf, err
}

func (f *FabricContractStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return f.stub.CreateCompositeKey(objectType, attributes)
}

func (f *FabricContractStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return f.stub.SplitCompositeKey(compositeKey)
}

func (f *FabricContractStub) GetTxTimestamp() (time.Time, error) {
	ts, err := f.stub.GetTxTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	if ts == nil {
		return time.Time{}, errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < -62135596800 {
		return time.Time{}, fmt.Errorf("timestamp: %v before 0001-01-01", ts)
	}
	if ts.Seconds >= 253402300800 {
		return time.Time{}, fmt.Errorf("timestamp: %v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return time.Time{}, fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", ts)
	}

	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC(), nil
}

func (f *FabricContractStub) SetEvent(name string, payload []byte) error {
	return f.stub.SetEvent(name, payload)
}
