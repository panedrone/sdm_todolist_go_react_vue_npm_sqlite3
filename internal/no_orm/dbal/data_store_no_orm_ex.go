package dbal

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var ds = &_DS{}

func (ds *_DS) initDb() (err error) {
	ds.db, err = sql.Open("sqlite3", "./db/todolist.sqlite")
	DS = func(ctx context.Context) DataStore {
		return ds
	}
	return
}

func GetInstance() DataStore {
	return ds
}

func OpenDB() error {
	return ds.Open()
}

func CloseDB() error {
	return ds.Close()
}
