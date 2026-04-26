package handlers

import (
	"errors"
	"net/http"
	"sdm_demo_todolist/internal/gorm/dbal"
	m "sdm_demo_todolist/internal/gorm/dbal/models"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type projectHandlers struct {
}

func NewGormProjectHandlers() etc.ProjectHandlers {
	return &projectHandlers{}
}

// ProjectCreate godoc
// @Summary Create a new project
// @Description Create a new project with the given name
// @Tags projects
// @Accept  json
// @Produce  json
// @Param project body request.Project true "Project Name"
// @Success 201
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects [post]
func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := dbal.CreateProject(ctx, &m.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// ProjectsReadAll godoc
// @Summary Get all projects
// @Description Get a list of all projects
// @Tags projects
// @Produce  json
// @Success 200 {array} models.Project
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects [get]
func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := dbal.ReadAll(ctx)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, all)
}

// ProjectRead godoc
// @Summary Get a project by ID
// @Description Get details of a project by its ID
// @Tags projects
// @Produce  json
// @Param p_id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects/{p_id} [get]
func (h *projectHandlers) ProjectRead(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	pr, err := dbal.ReadProject(ctx, uri.PId)
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

// ProjectUpdate godoc
// @Summary Update a project
// @Description Update an existing project's name by its ID
// @Tags projects
// @Accept  json
// @Produce  json
// @Param p_id path int true "Project ID"
// @Param project body request.Project true "Project Data"
// @Success 200
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects/{p_id} [put]
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
	if _, err := dbal.UpdateProject(ctx, pr); err != nil {
		resp.Abort500(ctx, err)
	}
}

// ProjectDelete godoc
// @Summary Delete a project
// @Description Delete a project by its ID
// @Tags projects
// @Param p_id path int true "Project ID"
// @Success 204
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects/{p_id} [delete]
func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := dbal.DeleteProject(ctx, &m.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
