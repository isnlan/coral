package lock

import (
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

func Mutex(id string) *redsync.Mutex {
	if _r == nil {
		return nil
	}
	return _r.NewMutex(id)
}
