package contract

import (
	"errors"
	"fmt"
	"strconv"
)

// AppName/Table<pk1,pk2,...>/ -> key
// AppName/Table_Count/
type Index struct {
	app   string
	name  string
	key   string
	table *Table
}

func NewIndex(app, name, key string) *Index {
	ck := NewTable(app, name, "prefix", "index")
	return &Index{app: app, name: name, key: key, table: ck}
}

func (index *Index) Save(stub IContractStub, prefix string, value []byte) (int, error) {
	count, err := index.Total(stub, prefix)
	if err != nil {
		return 0, err
	}

	idx := strconv.Itoa(count)
	currentCount := strconv.Itoa(count + 1)

	// Address_N : ID
	err = index.table.Insert(stub, []string{prefix, idx}, value)
	if err != nil {
		return 0, err
	}

	// Address: Count
	err = stub.PutState(index.makeCountKey(prefix), []byte(currentCount))
	return count + 1, err
}

func (index *Index) Update(stub IContractStub, prefix string, idx int, value []byte) error {
	if idx < 0 {
		return errors.New("index error")
	}

	count, err := index.Total(stub, prefix)
	if err != nil {
		return err
	}

	if idx > count-1 {
		return errors.New("out of range")
	}

	return index.table.Insert(stub, []string{prefix, strconv.Itoa(idx)}, value)
}

func (index *Index) Latest(stub IContractStub, prefix string) ([]byte, error) {
	count, err := index.Total(stub, prefix)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	idx := strconv.Itoa(count - 1)

	return index.table.GetValue(stub, []string{prefix, idx})
}

// idx start from 0
func (index *Index) GetByIndex(stub IContractStub, prefix string, idx int) ([]byte, error) {
	if idx < 0 {
		return nil, errors.New("index error")
	}

	count, err := index.Total(stub, prefix)
	if err != nil {
		return nil, err
	}

	if idx > count-1 {
		return nil, errors.New("out of range")
	}

	return index.table.GetValue(stub, []string{prefix, strconv.Itoa(idx)})
}

func (index *Index) Total(stub IContractStub, prefix string) (int, error) {
	key := index.makeCountKey(prefix)
	countBytes, err := stub.GetState(key)
	if err != nil {
		return 0, err
	}

	if countBytes == nil || len(countBytes) == 0 {
		return 0, nil
	}

	return strconv.Atoi(string(countBytes))
}

func (index *Index) List(stub IContractStub, prefix string, offset, limit int, order bool) ([][]byte, error) {
	count, err := index.Total(stub, prefix)
	if err != nil {
		return nil, err
	}

	var list [][]byte

	if count <= offset {
		return list, nil
	}

	j := 1
	if order {
		for i := offset; i < count; i++ {
			if limit > 0 && j > limit {
				break
			}
			value, err := index.getValue(stub, prefix, i)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
			j++
		}
	} else {
		for i := count - offset - 1; i >= 0; i-- {
			if limit > 0 && j > limit {
				break
			}
			value, err := index.getValue(stub, prefix, i)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
			j++
		}
	}

	return list, nil
}

func (index *Index) Filter(stub IContractStub, prefix string, order bool, f func(value []byte) (bool, error)) error {
	count, err := index.Total(stub, prefix)
	if err != nil {
		return err
	}

	if order {
		for i := 0; i < count; i++ {
			value, err := index.getValue(stub, prefix, i)
			if err != nil {
				return err
			}
			ctiu, err := f(value)
			if err != nil {
				return err
			}
			if !ctiu {
				return nil
			}
		}
	} else {
		for i := count - 1; i >= 0; i-- {
			value, err := index.getValue(stub, prefix, i)
			if err != nil {
				return err
			}
			ctiu, err := f(value)
			if err != nil {
				return err
			}
			if !ctiu {
				return nil
			}
		}
	}

	return nil
}

func (index *Index) getValue(stub IContractStub, prefix string, idx int) ([]byte, error) {
	return index.table.GetValue(stub, []string{prefix, strconv.Itoa(idx)})
}

func (index *Index) makeCountKey(prefix string) string {
	return fmt.Sprintf("%s|%s<count>/%s_%s", index.app, index.name, prefix, index.key)
}
