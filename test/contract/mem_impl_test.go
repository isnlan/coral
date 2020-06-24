package contract

import (
	"fmt"
	"testing"

	"github.com/snlansky/coral/pkg/contract/mock"

	"github.com/snlansky/coral/pkg/contract"
	"github.com/stretchr/testify/assert"
)

func TestNewIndex(t *testing.T) {
	chain := mock.NewMemoryFactoryChain()
	stub := chain.NewStub("address1")

	index := contract.NewIndex("MyAPP", "user_table", "uid")
	addr1 := "user1_addr"
	total, err := index.Total(stub, addr1)
	assert.NoError(t, err)
	assert.Equal(t, total, 0)

	count, err := index.Save(stub, addr1, []byte("user1_id_0"))
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	total, err = index.Total(stub, addr1)
	assert.NoError(t, err)
	assert.Equal(t, total, 1)

	{
		addr2 := "user2_addr"
		for i := 0; i < 10; i++ {
			count, err := index.Save(stub, addr2, []byte(fmt.Sprintf("user2_id_%d", i)))
			assert.NoError(t, err)
			assert.Equal(t, count, i+1)
		}
	}

	i, err := index.Total(stub, "user2_addr")
	assert.NoError(t, err)
	assert.Equal(t, i, 10)

	list, err := index.List(stub, "user2_addr", 0, 0, true)
	assert.NoError(t, err)
	for _, v := range list {
		fmt.Println(string(v))
	}

	chain.Debug()
}

func TestIndex_Update(t *testing.T) {
	chain := mock.NewMemoryFactoryChain()
	stub := chain.NewStub("address1")

	index := contract.NewIndex("MyAPP", "user_table", "uid")

	addr := "user1_addr"
	{
		for i := 0; i < 10; i++ {
			count, err := index.Save(stub, addr, []byte(fmt.Sprintf("user_id_%d", i)))
			assert.NoError(t, err)
			assert.Equal(t, count, i+1)
		}
	}

	{
		for i := 0; i < 10; i++ {
			err := index.Update(stub, addr, i, []byte(fmt.Sprintf("persion_id_%d", i)))
			assert.NoError(t, err)
		}
	}

	chain.Debug()
}

func TestNewMemImpl(t *testing.T) {
	chain := mock.NewMemoryFactoryChain()
	stub := chain.NewStub("address1")

	index := contract.NewIndex("MyAPP", "user_table", "uid")
	addr1 := "user1_addr"

	for i := 0; i < 100; i++ {
		_, err := index.Save(stub, addr1, []byte(fmt.Sprintf("user1_id_%d", i)))
		assert.NoError(t, err)
	}

	err := index.Filter(stub, addr1, false, func(value []byte) (bool, error) {
		fmt.Println(string(value))
		return true, nil
	})
	assert.NoError(t, err)
}

func TestNewMemImpl2(t *testing.T) {
	chain := mock.NewMemoryFactoryChain()
	stub := chain.NewStub("address1")

	index := contract.NewIndex("MyAPP", "user_table", "uid")

	{
		addr2 := "user2_addr"
		for i := 0; i < 10; i++ {
			count, err := index.Save(stub, addr2, []byte(fmt.Sprintf("user2_id_%d", i)))
			assert.NoError(t, err)
			assert.Equal(t, count, i+1)
		}
	}

	i, err := index.Total(stub, "user2_addr")
	assert.NoError(t, err)
	assert.Equal(t, i, 10)

	list, err := index.List(stub, "user2_addr", 2, 2, true)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 2)

	for _, v := range list {
		fmt.Println("value => ", string(v))
	}

	chain.Debug()
}

func TestNewMemImpl_Index_Latest(t *testing.T) {
	chain := mock.NewMemoryFactoryChain()
	stub := chain.NewStub("address1")

	index := contract.NewIndex("MyAPP", "user_table", "uid")

	{
		addr1 := "user1_addr"
		bytes, err := index.Latest(stub, addr1)
		assert.NoError(t, err)
		assert.True(t, len(bytes) == 0)
		assert.True(t, bytes == nil)
	}

	{
		addr2 := "user2_addr"
		for i := 0; i < 10; i++ {
			count, err := index.Save(stub, addr2, []byte(fmt.Sprintf("user2_id_%d", i)))
			assert.NoError(t, err)
			assert.Equal(t, count, i+1)
		}

		bytes, err := index.Latest(stub, addr2)
		assert.NoError(t, err)
		assert.Equal(t, bytes, []byte("user2_id_9"))
	}

	{
		addr2 := "user2_addr"
		bytes, err := index.GetByIndex(stub, addr2, 10)
		assert.Error(t, err)
		assert.True(t, len(bytes) == 0)

		bytes1, err1 := index.GetByIndex(stub, addr2, 0)
		assert.NoError(t, err1)
		assert.Equal(t, bytes1, []byte("user2_id_0"))

		bytes, err = index.GetByIndex(stub, addr2, 4)
		assert.NoError(t, err)
		assert.Equal(t, bytes, []byte("user2_id_4"))

		bytes, err = index.GetByIndex(stub, addr2, 9)
		assert.NoError(t, err)
		assert.Equal(t, bytes, []byte("user2_id_9"))
	}

	chain.Debug("MyAPP|user_table<prefix:index")
}
