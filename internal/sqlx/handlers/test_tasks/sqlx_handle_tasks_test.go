package test_tasks

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sdm_demo_todolist/internal/sqlx/handlers"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"sdm_demo_todolist/internal/sqlx/dbal"
	"sdm_demo_todolist/internal/sqlx/dbal/dto"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	origCreate := dbal.CreateTask
	origRead := dbal.ReadTask
	origUpdate := dbal.UpdateTask
	origDelete := dbal.DeleteTask
	origGetByProject := dbal.GetGroupTasks

	defer func() {
		dbal.CreateTask = origCreate
		dbal.ReadTask = origRead
		dbal.UpdateTask = origUpdate
		dbal.DeleteTask = origDelete
		dbal.GetGroupTasks = origGetByProject
	}()

	dbal.CreateTask = func(ctx context.Context, t *dto.Task) error {
		t.TID = 1
		return nil
	}

	dbal.ReadTask = func(ctx context.Context, id int64) (*dto.Task, error) {
		return &dto.Task{
			TID:       id,
			PID:       1,
			TSubject:  "demo",
			TPriority: 1,
			TDate:     "2024-01-01 10:00:00",
		}, nil
	}

	dbal.UpdateTask = func(ctx context.Context, t *dto.Task) (int64, error) {
		return 1, nil
	}

	dbal.DeleteTask = func(ctx context.Context, t *dto.Task) (int64, error) {
		return 1, nil
	}

	dbal.GetGroupTasks = func(ctx context.Context, pid int64) ([]*dto.TaskLi, error) {
		return []*dto.TaskLi{
			{TID: 1, TSubject: "t1"},
			{TID: 2, TSubject: "t2"},
		}, nil
	}

	h := handlers.NewSqlxTaskHandlers()

	router = gin.New()
	router.POST("/projects/:p_id/tasks", h.TaskCreate)
	router.GET("/projects/:p_id/tasks", h.TasksReadByProject)
	router.GET("/tasks/:t_id", h.TaskRead)
	router.PUT("/tasks/:t_id", h.TaskUpdate)
	router.DELETE("/tasks/:t_id", h.TaskDelete)

	code := m.Run()
	os.Exit(code)
}

func TestTaskCreate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/projects/1/tasks",
		strings.NewReader(`{"t_subject":"test task"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTasksReadByProject_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/projects/1/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "t1")
}

func TestTaskRead_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "demo")
}

func TestTaskRead_NotFound(t *testing.T) {
	orig := dbal.ReadTask
	defer func() { dbal.ReadTask = orig }()

	dbal.ReadTask = func(ctx context.Context, id int64) (*dto.Task, error) {
		return nil, sql.ErrNoRows
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTaskUpdate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/1",
		strings.NewReader(`{
			"p_id":1,
			"t_subject":"updated",
			"t_priority":2,
			"t_date":"2024-01-01 10:00:00"
		}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTaskUpdate_ValidationError(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/1",
		strings.NewReader(`{
			"p_id":1,
			"t_subject":"",
			"t_priority":0,
			"t_date":"bad-date"
		}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTaskDelete_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTaskCreate_DBError(t *testing.T) {
	orig := dbal.CreateTask
	defer func() { dbal.CreateTask = orig }()

	dbal.CreateTask = func(ctx context.Context, t *dto.Task) error {
		return errors.New("db error")
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects/1/tasks",
		strings.NewReader(`{"t_subject":"fail"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
