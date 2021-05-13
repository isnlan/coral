package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	Init("127.0.0.1:6379", "", 0)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		mutex := New(context.Background(), "abc")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock()
		if err != nil {
			panic(err)
		}

		fmt.Println("abc lock 1")
		time.Sleep(time.Second * 10)
		fmt.Println("abc release 1")
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex := New(context.Background(), "abc")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock()
		if err != nil {
			panic(err)
		}

		fmt.Println("abc lock 2")
		time.Sleep(time.Second * 10)
		fmt.Println("abc release 2")
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex := New(context.Background(), "xyz")
		if mutex == nil {
			fmt.Println("get mutex")
			return
		}

		err := mutex.Lock()
		if err != nil {
			panic(err)
		}

		fmt.Println("xyz lock 1")
		time.Sleep(time.Second * 10)
		fmt.Println("xyz release 1")
		mutex.Unlock()
	}()

	wg.Wait()
}
