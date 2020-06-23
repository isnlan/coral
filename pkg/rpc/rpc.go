package rpc

import "encoding/json"

type Rpc interface {
	Register(rcvr interface{}) error
	RegisterName(name string, rcvr interface{}) error
	Handler(req *Request, baseParam ...interface{}) (interface{}, error)
}

type Request struct {
	ServiceMethod string             `json:"func_name"` // format: "Service.Method"
	Params        []*json.RawMessage `json:"params"`
}

type ClientRequest struct {
	Method string        `json:"func_name"`
	Params []interface{} `json:"params"`
}
