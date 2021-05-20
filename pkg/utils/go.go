package utils

import (
	"encoding/json"
	"reflect"
)

func SafeRun(f func()) {
	defer func() {
		if re := recover(); re != nil {
			logger.Errorf("panic: %v", re)
		}
	}()
	f()
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
