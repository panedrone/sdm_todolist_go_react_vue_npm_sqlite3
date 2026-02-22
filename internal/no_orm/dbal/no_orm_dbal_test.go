package dbal

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var ctx context.Context

func TestMain(m *testing.M) {
	var err error
	ds.db, err = sql.Open("sqlite3", "../../todolist.sqlite")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = CloseDB()
	}()
	ctx = context.Background()
	code := m.Run()
	// .................... clean up
	os.Exit(code)
}

func Test_no_orm_ReadAllScalars(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := ds.query(ctx, sqlStr)
	if err != nil {
		return
	}
	var res []int64
	// Scan error on column index 0, name "p_id": destination not a pointer; sql: Scan error on column index 0,
	// name "p_id": destination not a pointer; sql: Scan error on column index 0, name "p_id": destination not a pointer
	err = ds.selectAllScalars(rows, &res)
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	for _, v := range res {
		if v < 1 {
			t.Fatalf("Error: %v", v)
			return
		}
	}
	t.Logf("%v\n", res)
}

func Test_no_orm_ReadAllScalarPtr(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := ds.query(ctx, sqlStr)
	if err != nil {
		return
	}
	defer _close(rows)
	var res []*int64
	err = ds.selectAllScalars(rows, &res)
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	var values []int64
	for _, v := range res {
		if *v < 1 {
			t.Fatalf("Error: %v", *v)
			return
		}
		values = append(values, *v)
	}
	t.Logf("%v\n", values)
}
