package xgin

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	ses "github.com/gin-contrib/sessions/redis"
)

func NewStore(uri string, db int, password string) (sessions.Store, error) {
	return ses.NewStoreWithDB(32, "tcp", uri, password, strconv.Itoa(db))
}
