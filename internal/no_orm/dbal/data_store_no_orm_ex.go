package dbal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var ds = &_DS{}

func (ds *_DS) initDb() (err error) {
	ds.db, err = sql.Open("sqlite3", "./db/todolist.sqlite")
	return
}

func OpenDB() error {
	return ds.Open()
}

func CloseDB() error {
	return ds.Close()
}

func NewProjectsDao() *ProjectsDao {
	return &ProjectsDao{ds: ds}
}

func NewTasksDao() *TasksDao {
	return &TasksDao{ds: ds}
}
