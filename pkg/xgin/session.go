package xgin

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
)

func NewRedisStore(uri string, db int, password string) (sessions.Store, error) {
	return redis.NewStoreWithDB(32, "tcp", uri, password, strconv.Itoa(db), []byte("secret"))
}

func NewMemStore() sessions.Store {
	return memstore.NewStore([]byte("secret"))
}
