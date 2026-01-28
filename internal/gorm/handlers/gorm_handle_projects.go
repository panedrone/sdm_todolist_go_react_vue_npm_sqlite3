package handlers

import (
	"errors"
	"net/http"
	dbal2 "sdm_demo_todolist/internal/gorm/dbal"
	m "sdm_demo_todolist/internal/gorm/dbal/models"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type projectHandlers struct {
	dao *dbal2.ProjectsDao
}

func NewGormProjectHandlers() etc.ProjectHandlers {
	return &projectHandlers{
		dao: dbal2.NewProjectsDao(),
	}
}

func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := h.dao.CreateProject(ctx, &m.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := h.dao.ReadAll(ctx)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, all)
}

func (h *projectHandlers) ProjectRead(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	pr, err := h.dao.ReadProject(ctx, uri.PId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Abort404NotFound(ctx, err)
			return
		}
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, pr)
}

func (h *projectHandlers) ProjectUpdate(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	pr := &m.Project{PID: uri.PId, PName: req.PName}
	if _, err := h.dao.UpdateProject(ctx, pr); err != nil {
		resp.Abort500(ctx, err)
	}
}

func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := h.dao.DeleteProject(ctx, &m.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
