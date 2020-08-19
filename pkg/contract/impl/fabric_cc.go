package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/snlansky/coral/pkg/contract"
	"github.com/snlansky/coral/pkg/rpc"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type FabricChaincode struct {
	rpc rpc.Rpc
}

func NewFabricChaincode() *FabricChaincode {
	return &FabricChaincode{rpc: rpc.New()}
}

func (cc *FabricChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("SUCCESS"))
}

func (cc *FabricChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	stb := NewFabricContractStub(stub)
	args := stb.GetArgs()
	if len(args) <= 0 || len(args) > 2 {
		return shim.Error(contract.ERR_PARAM_INVALID)
	}

	method := string(args[0])
	var param []*json.RawMessage

	if len(args) == 2 {
		err := json.Unmarshal(args[1], &param)
		if err != nil {
			log.Printf("ERR: json.Unmarshal error:%s, date:%s\n", err.Error(), string(args[1]))
			return shim.Error(contract.ERR_JSON_UNMARSHAL)
		}
	}

	addr, err := stb.GetAddress()
	if err != nil {
		log.Printf("ERR: auth user failed, error:%s\n", err.Error())
		return shim.Error("ERR_INVALID_CERT")
	}

	log.Printf("INFO: address:%s, method:%s, params:%v\n", addr, method, param)

	req := &rpc.Request{
		ServiceMethod: method,
		Params:        param,
	}

	return cc.handler(stb, req)
}

func (cc *FabricChaincode) handler(stub contract.IContractStub, req *rpc.Request) pb.Response {
	var (
		ret interface{}
		err error
	)

	startTime := time.Now()

	ret, err = cc.recoverHandler(stub, req)
	if err != nil {
		log.Printf("ERR:response error:%s\n", err.Error())
		return shim.Error(err.Error())
	}

	if ret == nil {
		log.Printf("INFO: process takes %v, response success:null\n", time.Since(startTime))
		return shim.Success(nil)
	}

	buf, err := json.Marshal(ret)
	if err != nil {
		log.Printf("ERR:response error:%s\n", err.Error())
		return shim.Error(contract.ERR_JSON_MARSHAL)
	}
	log.Printf("INFO: process takes %v, response success:%s\n", time.Since(startTime), string(buf))
	return shim.Success(buf)
}

func (cc *FabricChaincode) recoverHandler(stub contract.IContractStub, req *rpc.Request) (ret interface{}, err error) {
	defer func() {
		if re := recover(); re != nil {
			switch v := re.(type) {
			case contract.InternalError:
				//logger.Error(v.Error())
				err = errors.New(v.External())
			case *contract.InternalError:
				//logger.Error(v.Error())
				err = errors.New(v.External())
			case string:
				err = errors.New(v)
			case runtime.Error:
				fmt.Printf("ERR: runtime error:%v \nstack :%s\n", v, string(debug.Stack()))
				err = contract.ErrRuntime
			case error:
				err = v
			default:
				err = fmt.Errorf("ERR: ohter error type:%v, value:%v", reflect.TypeOf(re), v)
			}
		}
	}()

	ret, err = cc.rpc.Handler(req, stub)
	return
}

func (cc *FabricChaincode) Register(i interface{}) {
	err := cc.rpc.Register(i)
	if err != nil {
		panic(err)
	}
}

func (cc *FabricChaincode) Start() {
	err := shim.Start(cc)
	if err != nil {
		panic("Error starting chaincode - " + err.Error())
	}
}
