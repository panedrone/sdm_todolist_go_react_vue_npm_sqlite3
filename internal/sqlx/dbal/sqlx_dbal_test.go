package dbal

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var ctx context.Context

func TestMain(m *testing.M) {

	db, err := sqlx.Open("sqlite3", "../../todolist.sqlite")
	if err != nil {
		println(err.Error())
		return
	}
	ds = &_DS{db: db}
	defer func() {
		_ = db.Close()
	}()

	ctx = context.Background()
	code := m.Run()
	// .................... clean up
	os.Exit(code)
}

func Test_sqlx_ReadAllScalars(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := ds.querySqlx(ctx, sqlStr)
	if err != nil {
		return
	}
	defer _close(rows)
	var res []int64
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

func Test_sqlx_ReadAllScalarPtr(t *testing.T) {
	sqlStr := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	rows, err := ds.querySqlx(ctx, sqlStr)
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

//func Test_sqlx_GetProjectIds(t *testing.T) {
//	dao := NewProjectsDao()
//	res, err := dao.GetProjectIds(ctx)
//	if err != nil {
//		t.Fatalf("Error: %s", err)
//		return
//	}
//	for _, v := range res {
//		if v < 1 {
//			t.Fatalf("Error: %v", v)
//			return
//		}
//	}
//	t.Logf("%v\n", res)
//}
