package errors

import "fmt"

func Check(err error, msg ...interface{}) {
	if err != nil {
		var info string
		var code int
		if e, ok := err.(CodeError); ok {
			info = e.Error()
			code = e.Code()
		} else {
			info = err.Error()
			code = InternalErrorCode
		}

		if msg != nil && len(msg) > 0 {
			info = fmt.Sprintf("%s, %v", info, msg)
		}
		panic(NewWithInfo(code, info))
	}
}

func Throw(desc string, code int, msg ...interface{}) {
	if msg != nil && len(msg) > 0 {
		desc = fmt.Sprintf("%s, %v", desc, msg)
	}
	panic(NewWithInfo(code, desc))
}
