package sqldb

import (
	"database/sql"
)

func New(driver, dsn string) *sql.DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
