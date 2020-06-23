package errors

func Check(err error, msg ...interface{}) {
	if err != nil {
		if e, ok := err.(CodeError); ok {
			panic(e)
		} else {
			panic(NewWithCode(err, InternalErrorCode))
		}

	}
}

func Throw(desc string, code int, msg ...interface{}) {
	panic(NewWithInfo(desc, code))
}
