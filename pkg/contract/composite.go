package contract

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	minUnicodeRuneValue   = 0            //U+0000
	maxUnicodeRuneValue   = utf8.MaxRune //U+10FFFF - maximum (and unallocated) code point
	compositeKeyNamespace = "\x00"
	emptyKeySubstitute    = "\x01"
	keySep                = "/"
)

func CreateCompositeKey(objectType string, attributes []string) (string, error) {
	if err := validateCompositeKeyAttribute(objectType); err != nil {
		return "", err
	}
	ck := compositeKeyNamespace + objectType + string(minUnicodeRuneValue)
	for _, att := range attributes {
		if err := validateCompositeKeyAttribute(att); err != nil {
			return "", err
		}
		ck += att + string(minUnicodeRuneValue)
	}
	return ck, nil
}

func SplitCompositeKey(compositeKey string) (string, []string, error) {
	componentIndex := 1
	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeRuneValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	return components[0], components[1:], nil
}

func validateCompositeKeyAttribute(str string) error {
	if !utf8.ValidString(str) {
		return errors.New(fmt.Sprintf("not a valid utf8 string: [%x]", str))
	}
	for index, runeValue := range str {
		if runeValue == minUnicodeRuneValue || runeValue == maxUnicodeRuneValue {
			return errors.New(fmt.Sprintf(`input contain unicode %#U starting at position [%d]. %#U and %#U are not allowed in the input attribute of a composite key`,
				runeValue, index, minUnicodeRuneValue, maxUnicodeRuneValue))
		}
	}
	return nil
}

func CreateKey(objectType string, attributes []string) (string, error) {
	var components []string
	components = append(components, objectType)
	components = append(components, attributes...)
	err := validateKeyAttribute(components)
	if err != nil {
		return "", err
	}
	return strings.Join(components, keySep), nil
}

func SplitKey(key string) (string, []string, error) {
	keys := strings.Split(key, keySep)
	return keys[0], keys[1:], nil
}

func validateKeyAttribute(attrs []string) error {
	for _, k := range attrs {
		if strings.Contains(k, keySep) {
			return fmt.Errorf("cannot contains the special delimiter '%s'", keySep)
		}
	}
	return nil
}
