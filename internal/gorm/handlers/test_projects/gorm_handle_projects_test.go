package test_projects

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sdm_demo_todolist/internal/gorm/dbal"
	"sdm_demo_todolist/internal/gorm/dbal/models"
	"sdm_demo_todolist/internal/gorm/handlers"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

	dbal.CreateProject = func(ctx context.Context, item *models.Project) error {
		item.PID = 1
		return nil
	}

	dbal.ReadAll = func(ctx context.Context) ([]*models.ProjectLi, error) {
		return []*models.ProjectLi{
			{PID: 1, PName: "p1", PTasksCount: 0},
			{PID: 2, PName: "p2", PTasksCount: 5},
		}, nil
	}

	dbal.ReadProject = func(ctx context.Context, id int64) (*models.Project, error) {
		return &models.Project{PID: id, PName: "demo"}, nil
	}

	dbal.UpdateProject = func(ctx context.Context, item *models.Project) (int64, error) {
		return 1, nil
	}

	dbal.DeleteProject = func(ctx context.Context, item *models.Project) (int64, error) {
		return 1, nil
	}

	h := handlers.NewGormProjectHandlers()

	router = gin.New()
	router.POST("/projects", h.ProjectCreate)
	router.GET("/projects", h.ProjectsReadAll)
	router.GET("/projects/:p_id", h.ProjectRead)
	router.PUT("/projects/:p_id", h.ProjectUpdate)
	router.DELETE("/projects/:p_id", h.ProjectDelete)

	code := m.Run()
	os.Exit(code)
}

func TestProjectCreate_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/projects",
		strings.NewReader(`{"p_name":"test project"}`),
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

	dbal.ReadProject = func(ctx context.Context, id int64) (*models.Project, error) {
		return nil, gorm.ErrRecordNotFound
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
