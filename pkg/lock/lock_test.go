package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/isnlan/coral/pkg/db"
)

func TestMutex(t *testing.T) {
	err := db.InitRedis("127.0.0.1:6379", 0, "")
	if err != nil {
		panic(err)
	}
	Init(db.GetClient())

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		mutex := New("abc")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Println("abc lock 1")
		time.Sleep(time.Second * 3)
		fmt.Println("abc release 1")
		mutex.Unlock(context.Background())
	}()

	go func() {
		defer wg.Done()
		mutex := New("abc")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Println("abc lock 2")
		time.Sleep(time.Second * 3)
		fmt.Println("abc release 2")
		mutex.Unlock(context.Background())
	}()

	go func() {
		defer wg.Done()
		mutex := New("xyz")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Println("xyz lock 1")
		time.Sleep(time.Second * 10)
		fmt.Println("xyz release 1")
		mutex.Unlock(context.Background())
	}()

	wg.Wait()
}
