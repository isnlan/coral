package contract

import (
	"errors"
	"strings"
)

const (
	ERR_RUNTIME                = "ERR_RUNTIME"                // 运行时错误
	ERR_INTERNAL_INVALID       = "ERR_INTERNAL_INVALID"       // 内部错误
	ERR_NOT_FIND_INIT_FUNCTION = "ERR_NOT_FIND_INIT_FUNCTION" // 没有没找初始化函数
	ERR_PARSE_RPC_REQ          = "ERR_PARSE_RPC_REQ"          // 解析RPC请求错误
	ERR_METHOD_NOT_FOUND       = "ERR_METHOD_NOT_FOUND"       // 没有找到方法
	ERR_PARAM_COUNT_NOT_MATCH  = "ERR_PARAM_COUNT_NOT_MATCH"  // 调用参数数量不匹配
	ERR_PARAM_INVALID          = "ERR_PARAM_INVALID"          // 参数错误
	ERR_JSON_MARSHAL           = "ERR_JSON_MARSHAL"           // json数据错误
	ERR_JSON_UNMARSHAL         = "ERR_JSON_UNMARSHAL"         // 读取json数据错误
)

var (
	ErrRuntime             = errors.New(ERR_RUNTIME)
	ErrInternalInvalid     = errors.New(ERR_INTERNAL_INVALID)
	ErrNotFindInitFunction = errors.New(ERR_NOT_FIND_INIT_FUNCTION)
	ErrParseRpcReq         = errors.New(ERR_PARSE_RPC_REQ)
	ErrMethodNotFound      = errors.New(ERR_METHOD_NOT_FOUND)
	ErrParamCountNotMatch  = errors.New(ERR_PARAM_COUNT_NOT_MATCH)
	ErrParamInvalid        = errors.New(ERR_PARAM_INVALID)
	ErrJsonMarshal         = errors.New(ERR_JSON_MARSHAL)
	ErrJsonUnmarshal       = errors.New(ERR_JSON_UNMARSHAL)
)

type InternalError struct {
	err  error
	info []string
}

func NewInternalError(err error, info ...string) *InternalError {
	return &InternalError{err: err, info: info}
}

func (e *InternalError) Error() string {
	msg := e.info[:]
	if e.err != nil {
		msg = append(msg, e.err.Error())
	}
	return strings.Join(msg, ",")
}

func (e *InternalError) External() string {
	return strings.Join(e.info, ",")
}

func Check(err error, info ...string) {
	if err != nil {
		switch e := err.(type) {
		case *InternalError:
			panic(e)
		default:
			panic(&InternalError{err: err, info: info})
		}
	}
}

func Throw(info ...string) {
	panic(NewInternalError(nil, info...))
}
