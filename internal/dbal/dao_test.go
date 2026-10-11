package dbal

// Все сгенерированные DAO на db/todolist.sqlite. Всё идёт в одной транзакции, которая
// откатывается в конце: база остаётся как была.

import (
	"context"
	"os"
	"testing"

	"sdm_demo_todolist/internal/dbal/models"
)

var ctx context.Context

func TestMain(m *testing.M) {
	_ = os.Chdir("../..") // ./db/todolist.sqlite — от корня проекта
	if err := Ds().Open(); err != nil {
		panic(err.Error())
	}
	var err error
	ctx, err = Ds().Begin(context.Background())
	if err != nil {
		panic(err.Error())
	}
	code := m.Run()
	_ = Ds().Rollback(&ctx)
	_ = Ds().Close()
	os.Exit(code)
}

// res — результат DAO-функции целиком: r(dbal.F(...)) передаётся в проверку одним аргументом
type res[T any] struct {
	v   T
	err error
}

func r[T any](v T, err error) res[T] { return res[T]{v, err} }

func must[T any](t *testing.T, what string, x res[T]) T {
	t.Helper()
	if x.err != nil {
		t.Fatalf("%s: %v", what, x.err)
	}
	return x.v
}

func noErr(t *testing.T, what string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
}

func one(t *testing.T, what string, x res[int64]) {
	t.Helper()
	if x.err != nil || x.v != 1 {
		t.Fatalf("%s: rows affected %d, %v", what, x.v, x.err)
	}
}

// Тесты зависят друг от друга по данным, поэтому — один сценарий по шагам.
func TestTodolist(t *testing.T) {
	// projects — с чистого листа (внутри транзакции): запросы скаляров ниже ждут одну строку.
	// Задачи остаются (среди них есть и без проекта), счёт — от base.
	for _, p := range must(t, "ReadAll", r(ReadAll(ctx))) {
		must(t, "DelByProject", r(DelByProject(ctx, p.PID)))
		one(t, "DeleteProject", r(DeleteProject(ctx, &models.Project{PID: p.PID})))
	}
	base := must(t, "GetCount", r(GetCount(ctx)))

	// projects: CRUD
	p := &models.Project{PName: "Test project"}
	noErr(t, "CreateProject", CreateProject(ctx, p))
	if p.PID == 0 {
		t.Fatal("CreateProject: p_id is not returned")
	}
	p.PName = "Renamed"
	one(t, "UpdateProject", r(UpdateProject(ctx, p)))
	if got := must(t, "ReadProject", r(ReadProject(ctx, p.PID))); *got != *p {
		t.Errorf("ReadProject: %+v, want %+v", *got, *p)
	}

	// одна строка в projects — запросы скаляров
	if id := must(t, "PID", r(PID(ctx))); id != p.PID {
		t.Errorf("PID: %d", id)
	}
	if id := must(t, "ProjectID", r(ProjectID(ctx, p.PID))); id != p.PID {
		t.Errorf("ProjectID: %d", id)
	}
	if id := must(t, "GetPID", r(GetPID(ctx, p.PID))); id != p.PID {
		t.Errorf("GetPID: %d", id)
	}
	if ids := must(t, "GetPIds", r(GetPIds(ctx, p.PID))); len(ids) != 1 || ids[0] != p.PID {
		t.Errorf("GetPIds: %v", ids)
	}
	for name, f := range map[string]func(context.Context) ([]int64, error){
		"GetProjectIds": GetProjectIds, "GetProjectID": GetProjectID,
	} {
		if ids := must(t, name, r(f(ctx))); len(ids) != 1 || ids[0] != p.PID {
			t.Errorf("%s: %v", name, ids)
		}
	}

	// tasks: CRUD
	task := &models.Task{PID: p.PID, TPriority: 2, TDate: "2026-10-11", TSubject: "Test task", TComments: "c"}
	noErr(t, "CreateTask", CreateTask(ctx, task))
	if task.TID == 0 {
		t.Fatal("CreateTask: t_id is not returned")
	}
	task.TPriority, task.TSubject = 5, "Updated task"
	one(t, "UpdateTask", r(UpdateTask(ctx, task)))
	if got := must(t, "ReadTask", r(ReadTask(ctx, task.TID))); *got != *task {
		t.Errorf("ReadTask: %+v, want %+v", *got, *task)
	}
	task2 := &models.Task{PID: p.PID, TPriority: 1, TDate: "2026-10-12", TSubject: "Second task"}
	noErr(t, "CreateTask", CreateTask(ctx, task2))

	if list := must(t, "ReadTaskList", r(ReadTaskList(ctx))); int64(len(list)) != base+2 {
		t.Errorf("ReadTaskList: %d tasks, base %d", len(list), base)
	}
	li := must(t, "ReadByProject", r(ReadByProject(ctx, p.PID)))
	if len(li) != 2 || li[0].TID != task.TID || li[0].TSubject != "Updated task" || li[1].TPriority != 1 {
		t.Errorf("ReadByProject: %+v", li)
	}
	if all := must(t, "ReadAll", r(ReadAll(ctx))); len(all) != 1 || all[0].PTasksCount != 2 || all[0].PName != "Renamed" {
		t.Errorf("ReadAll: %+v", all)
	}
	if n := must(t, "GetCount", r(GetCount(ctx))); n != base+2 {
		t.Errorf("GetCount: %d, base %d", n, base)
	}

	one(t, "DeleteTask", r(DeleteTask(ctx, task2)))
	if _, err := ReadTask(ctx, task2.TID); err == nil {
		t.Error("ReadTask after DeleteTask: error expected")
	}

	// "testing" в dao_tasks.xml: запросы на delete_project_tasks.sql — строк не возвращают,
	// но delete выполняют
	if list := must(t, "GetProjectTasks2", r(GetProjectTasks2(ctx, p.PID))); len(list) != 0 {
		t.Errorf("GetProjectTasks2: %+v", list)
	}
	if n := must(t, "GetCount", r(GetCount(ctx))); n != base {
		t.Errorf("GetCount after GetProjectTasks2: %d, base %d", n, base)
	}
	if _, err := GetTask(ctx, p.PID); err == nil {
		t.Error("GetTask: no rows, error expected")
	}

	noErr(t, "CreateTask", CreateTask(ctx, task))
	if n := must(t, "DelByProject", r(DelByProject(ctx, p.PID))); n != 1 {
		t.Errorf("DelByProject: %d", n)
	}
	one(t, "DeleteProject", r(DeleteProject(ctx, p)))
	if _, err := ReadProject(ctx, p.PID); err == nil {
		t.Error("ReadProject after DeleteProject: error expected")
	}
}
