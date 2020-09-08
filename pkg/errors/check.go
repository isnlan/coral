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
			info = e.Error()
			code = InternalErrorCode
		}

		info = fmt.Sprintf("%s, %v", info, msg)
		panic(NewWithInfo(info, code))
	}
}

func Throw(desc string, code int, msg ...interface{}) {
	panic(NewWithInfo(fmt.Sprintf("%s, %v", desc, msg), code))
}
