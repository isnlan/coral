package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func SafeRun(f func() error) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("recover from panic, data: %v", re)
		}
	}()
	return f()
}

func MakeTypeName(tpy interface{}) string {
	return reflect.TypeOf(tpy).Elem().Name()
}

func MustStructToMap(data interface{}) map[string]interface{} {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var ret map[string]interface{}
	if err := json.Unmarshal(bytes, &ret); err != nil {
		panic(err)
	}
	return ret
}
