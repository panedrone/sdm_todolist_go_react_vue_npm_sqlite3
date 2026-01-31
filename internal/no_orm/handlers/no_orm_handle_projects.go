package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"sdm_demo_todolist/internal/no_orm/dbal"
	"sdm_demo_todolist/internal/no_orm/dbal/dto"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
)

type projectHandlers struct {
}

func NewNoOrmProjectHandlers() etc.ProjectHandlers {
	return &projectHandlers{}
}

func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := dbal.CreateProject(ctx, &dto.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := dbal.ReadAll(ctx)
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
	pr, err := dbal.ReadProject(ctx, uri.PId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	pr := &dto.Project{PID: uri.PId, PName: req.PName}
	if _, err := dbal.UpdateProject(ctx, pr); err != nil {
		resp.Abort500(ctx, err)
	}
}

func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := dbal.DeleteProject(ctx, &dto.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
