package contract

import "strings"

func CreateEvent(stub IContractStub, appName, eventName string, payload []byte) error {
	return stub.SetEvent(strings.Join([]string{appName, eventName}, "."), payload)
}

func MakeEventName(appName, eventName string) string {
	return strings.Join([]string{appName, eventName}, ".")
}
