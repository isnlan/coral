package contract

import (
	"fmt"
	"strings"
)

type KV struct {
	Keys  []string
	Value []byte
}

// AppName/Table<pk1,pk2,...>/pkv1:pkv2 -> value
type Table struct {
	app    string
	table  string
	fields []string
}

func NewTable(app string, table string, fields ...string) *Table {
	return &Table{app: app, table: table, fields: fields}
}

func (t *Table) GetType() string {
	return fmt.Sprintf("%s|%s<%s>", t.app, t.table, strings.Join(t.fields, ":"))
}

func (t *Table) createCompositeKey(stub IContractStub, keys []string) (string, error) {
	if err := t.check(keys); err != nil {
		return "", err
	}
	return stub.CreateCompositeKey(t.GetType(), keys)
}

func (t *Table) Insert(stub IContractStub, keys []string, value []byte) error {
	key, err := t.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	if value == nil {
		value = []byte{0x00}
	}
	return stub.PutState(key, value)
}

func (t *Table) Delete(stub IContractStub, keys []string) error {
	key, err := t.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	_, err = stub.DelState(key)
	return err
}

func (t *Table) Update(stub IContractStub, keys []string, value []byte) error {
	key, err := t.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	return stub.PutState(key, value)
}

func (t *Table) SplitKey(stub IContractStub, key string) ([]string, error) {
	_, ks, err := stub.SplitCompositeKey(key)
	if err != nil {
		return nil, err
	}
	err = t.check(ks)
	if err != nil {
		return nil, err
	}
	return ks, nil
}

func (t *Table) GetValue(stub IContractStub, keys []string) ([]byte, error) {
	key, err := t.createCompositeKey(stub, keys)
	if err != nil {
		return nil, err
	}
	return stub.GetState(key)
}

func (t *Table) check(keys []string) error {
	if len(keys) != len(t.fields) {
		return fmt.Errorf("keys count not matched. got %d, need %d", len(keys), len(t.fields))
	}
	return nil
}
