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
	mu       *redsync.Mutex
	extend   time.Duration
	releaseC chan struct{}
}

func New(id string) *Mutex {
	if _r == nil {
		return nil
	}

	extend := time.Second * 4
	mu := _r.NewMutex(id, redsync.SetExpiry(extend*2))
	return &Mutex{mu: mu, extend: extend, releaseC: make(chan struct{})}
}

func (m *Mutex) Lock(ctx context.Context) error {
	err := m.mu.Lock()
	if err != nil {
		return err
	}

	go func() {
		t := time.NewTicker(m.extend)
		for {
			select {
			case <-t.C:
				extend, err := m.mu.Extend()
				if err != nil {
					return
				}
				if !extend {
					return
				}
			case <-ctx.Done():
				return
			case <-m.releaseC:
				return
			}
		}
	}()
	return nil
}

func (m *Mutex) Unlock() (bool, error) {
	close(m.releaseC)
	return m.mu.Unlock()
}
