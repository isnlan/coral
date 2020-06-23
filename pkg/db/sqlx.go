package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _sqlxDB *sqlx.DB

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func InitSqlxDB(src string) *sqlx.DB {
	db := sqlx.MustConnect("mysql", src)
	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(2)
	// https://github.com/go-sql-driver/mysql/issues/446
	db.SetConnMaxLifetime(time.Second * 14400)
	return db
}

func GetSqlxDB() *sqlx.DB {
	return _sqlxDB
}
