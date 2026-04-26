package test_projects

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"sdm_demo_todolist/internal/no_orm/dbal"
	"sdm_demo_todolist/internal/no_orm/dbal/dto"
	"sdm_demo_todolist/internal/no_orm/handlers"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	origCreate := dbal.CreateProject
	origReadAll := dbal.ReadAll
	origRead := dbal.ReadProject
	origUpdate := dbal.UpdateProject
	origDelete := dbal.DeleteProject

	defer func() {
		dbal.CreateProject = origCreate
		dbal.ReadAll = origReadAll
		dbal.ReadProject = origRead
		dbal.UpdateProject = origUpdate
		dbal.DeleteProject = origDelete
	}()

	dbal.CreateProject = func(ctx context.Context, item *dto.Project) error {
		item.PID = 1
		return nil
	}

	dbal.ReadAll = func(ctx context.Context) ([]*dto.ProjectLi, error) {
		return []*dto.ProjectLi{
			{PID: 1, PName: "p1", PTasksCount: 0},
			{PID: 2, PName: "p2", PTasksCount: 5},
		}, nil
	}

	dbal.ReadProject = func(ctx context.Context, id int64) (*dto.Project, error) {
		return &dto.Project{PID: id, PName: "demo"}, nil
	}

	dbal.UpdateProject = func(ctx context.Context, item *dto.Project) (int64, error) {
		return 1, nil
	}

	dbal.DeleteProject = func(ctx context.Context, item *dto.Project) (int64, error) {
		return 1, nil
	}

	h := handlers.NewNoOrmProjectHandlers()

	router = gin.New()
	router.POST("/no_orm/projects", h.ProjectCreate)
	router.GET("/no_orm/projects", h.ProjectsReadAll)
	router.GET("/no_orm/projects/:p_id", h.ProjectRead)
	router.PUT("/no_orm/projects/:p_id", h.ProjectUpdate)
	router.DELETE("/no_orm/projects/:p_id", h.ProjectDelete)

	code := m.Run()
	os.Exit(code)
}

func TestProjectCreate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/no_orm/projects",
		strings.NewReader(`{"p_name":"test project"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProjectsReadAll_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/no_orm/projects", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "p1")
}

func TestProjectRead_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/no_orm/projects/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "demo")
}

func TestProjectRead_NotFound(t *testing.T) {
	orig := dbal.ReadProject
	defer func() { dbal.ReadProject = orig }()

	dbal.ReadProject = func(ctx context.Context, id int64) (*dto.Project, error) {
		return nil, sql.ErrNoRows
	}

	req := httptest.NewRequest(http.MethodGet, "/no_orm/projects/42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectUpdate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPut,
		"/no_orm/projects/1",
		strings.NewReader(`{"p_name":"updated"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectDelete_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/no_orm/projects/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
