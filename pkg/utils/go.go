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

func elemName(elem reflect.Type) string {
	if elem.Kind() == reflect.Ptr {
		return elem.Elem().Name()
	}

	return elem.Name()
}

func MakeTypeName(tpy interface{}) string {
	tp := reflect.TypeOf(tpy)

	switch tp.Kind() {
	case reflect.Map:
		return fmt.Sprintf("Map<%s:%s>", elemName(tp.Key()), elemName(tp.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("List<%s>", elemName(tp.Elem()))
	default:
		return tp.Elem().Name()
	}
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
