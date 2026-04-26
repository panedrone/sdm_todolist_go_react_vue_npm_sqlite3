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

// ProjectCreate godoc
// @Summary Create a new project (No-ORM)
// @Description Create a new project using raw SQL
// @Tags projects-no-orm
// @Accept  json
// @Produce  json
// @Param project body request.Project true "Project Name"
// @Success 201
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects [post]
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

// ProjectsReadAll godoc
// @Summary Get all projects (No-ORM)
// @Description Get a list of all projects using raw SQL
// @Tags projects-no-orm
// @Produce  json
// @Success 200 {array} dto.ProjectLi
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects [get]
func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := dbal.ReadAll(ctx)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, all)
}

// ProjectRead godoc
// @Summary Get a project by ID (No-ORM)
// @Description Get details of a project by its ID using raw SQL
// @Tags projects-no-orm
// @Produce  json
// @Param p_id path int true "Project ID"
// @Success 200 {object} dto.Project
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects/{p_id} [get]
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

// ProjectUpdate godoc
// @Summary Update a project (No-ORM)
// @Description Update an existing project's name by its ID using raw SQL
// @Tags projects-no-orm
// @Accept  json
// @Produce  json
// @Param p_id path int true "Project ID"
// @Param project body request.Project true "Project Data"
// @Success 200
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects/{p_id} [put]
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

// ProjectDelete godoc
// @Summary Delete a project (No-ORM)
// @Description Delete a project by its ID using raw SQL
// @Tags projects-no-orm
// @Param p_id path int true "Project ID"
// @Success 204
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects/{p_id} [delete]
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
