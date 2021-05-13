package lock

import (
	"context"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

var _r *redsync.Redsync

func Init(addr string, password string, db int) {
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	_r = redsync.New([]redsync.Pool{pool})
}

type Mutex struct {
	mu *redsync.Mutex
}

func New(ctx context.Context, id string) *Mutex {
	if _r == nil {
		return nil
	}

	mu := _r.NewMutex(id, redsync.SetExpiry(time.Second*8))

	go func() {
		t := time.NewTicker(time.Second * 4)
		for {
			select {
			case <-t.C:
				extend, err := mu.Extend()
				if err != nil {
					return
				}
				if !extend {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return &Mutex{mu: mu}
}

func (m *Mutex) Lock() error {
	return m.mu.Lock()
}

func (m *Mutex) Unlock() (bool, error) {
	return m.mu.Unlock()
}
