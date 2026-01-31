package test_projects

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

var (
	router *gin.Engine
)

// ---------- TestMain ----------

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// save originals
	origCreate := dbal.CreateProject
	origRead := dbal.ReadProject
	origUpdate := dbal.UpdateProject
	origDelete := dbal.DeleteProject
	origGet := dbal.GetProjects

	// restore after all test_tasks
	defer func() {
		dbal.CreateProject = origCreate
		dbal.ReadProject = origRead
		dbal.UpdateProject = origUpdate
		dbal.DeleteProject = origDelete
		dbal.GetProjects = origGet
	}()

	// default mocks (можно переопределять в тестах)
	dbal.CreateProject = func(ctx context.Context, p *dto.Project) error {
		p.PID = 1
		return nil
	}

	dbal.ReadProject = func(ctx context.Context, id int64) (*dto.Project, error) {
		return &dto.Project{PID: id, PName: "demo"}, nil
	}

	dbal.UpdateProject = func(ctx context.Context, p *dto.Project) (int64, error) {
		return 1, nil
	}

	dbal.DeleteProject = func(ctx context.Context, p *dto.Project) (int64, error) {
		return 1, nil
	}

	dbal.GetProjects = func(ctx context.Context) ([]*dto.ProjectLi, error) {
		return []*dto.ProjectLi{
			{PID: 1, PName: "p1"},
			{PID: 2, PName: "p2"},
		}, nil
	}

	// router + handlers
	h := handlers.NewSqlxProjectHandlers()

	router = gin.New()
	router.POST("/projects", h.ProjectCreate)
	router.GET("/projects", h.ProjectsReadAll)
	router.GET("/projects/:p_id", h.ProjectRead)
	router.PUT("/projects/:p_id", h.ProjectUpdate)
	router.DELETE("/projects/:p_id", h.ProjectDelete)

	code := m.Run()
	os.Exit(code)
}

//
// ---------- TESTS ----------
//

func TestProjectCreate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/projects",
		strings.NewReader(`{"p_name":"test"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProjectsReadAll_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "p1")
}

func TestProjectRead_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/projects/42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectUpdate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPut,
		"/projects/1",
		strings.NewReader(`{"p_name":"updated"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectDelete_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/projects/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProjectCreate_DBError(t *testing.T) {
	orig := dbal.CreateProject
	defer func() { dbal.CreateProject = orig }()

	dbal.CreateProject = func(ctx context.Context, p *dto.Project) error {
		return errors.New("db down")
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects",
		strings.NewReader(`{"p_name":"fail"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
