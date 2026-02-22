package dbal

import (
	"context"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ctx context.Context

func TestMain(m *testing.M) {
	var err error
	_ds.rootDb, err = gorm.Open(sqlite.Open("../../todolist.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
	ctx = context.Background()
	code := m.Run()
	// .................... clean up
	_ = CloseDB()
	os.Exit(code)
}

func Test_gorm_ReadAllScalars(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := _ds.rawQuery(ctx, sqlStr)
	if err != nil {
		return
	}
	var res []int64
	// Scan error on column index 0, name "p_id": destination not a pointer; sql: Scan error on column index 0,
	// name "p_id": destination not a pointer; sql: Scan error on column index 0, name "p_id": destination not a pointer
	err = _ds.selectAllScalars(rows, &res)
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

func Test_gorm_ReadAllScalarPtr(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := _ds.rawQuery(ctx, sqlStr)
	if err != nil {
		return
	}
	defer _close(rows)
	var res []*int64
	err = _ds.selectAllScalars(rows, &res)
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
