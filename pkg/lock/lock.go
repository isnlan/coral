package lock

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	redis "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
)

var _r *redsync.Redsync

func Init(client *redis.Client) {
	pool := goredis.NewPool(client)
	_r = redsync.New(pool)
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
	mu := _r.NewMutex(id, redsync.WithExpiry(extend*2))
	return &Mutex{mu: mu, extend: extend, releaseC: make(chan struct{})}
}

func (m *Mutex) Lock(ctx context.Context) error {
	err := m.mu.LockContext(ctx)
	if err != nil {
		return err
	}

	go func() {
		t := time.NewTicker(m.extend)
		for {
			select {
			case <-t.C:
				extend, err := m.mu.ExtendContext(ctx)
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

func (m *Mutex) Unlock(ctx context.Context) (bool, error) {
	close(m.releaseC)
	return m.mu.UnlockContext(ctx)
}
