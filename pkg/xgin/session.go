package xgin

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func UseSession(router *gin.Engine, store sessions.Store) {
	if store == nil {
		store = memstore.NewStore([]byte("secret"))
	}
	//store.Options(sessions.Options{
	//	Path:     "/",
	//	Domain:   "",
	//	MaxAge:   86400 * 30,
	//	Secure:   true,
	//	HttpOnly: true,
	//	SameSite: http.SameSiteNoneMode,
	//})
	router.Use(sessions.Sessions("session", store))
}

func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Default(ctx)
	session.Set(k, o)
	_ = session.Save()
}

func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Default(ctx)
	return session.Get(k)
}

func ClearAllSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	_ = session.Save()
}

func NewRedisStore(uri string, db int, password string) (sessions.Store, error) {
	return redis.NewStoreWithDB(32, "tcp", uri, password, strconv.Itoa(db), []byte("secret"))
}

func NewMemStore() sessions.Store {
	return memstore.NewStore([]byte("secret"))
}
