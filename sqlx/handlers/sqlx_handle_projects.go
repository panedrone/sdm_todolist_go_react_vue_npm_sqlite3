package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"sdm_demo_todolist/shared"
	"sdm_demo_todolist/shared/request"
	"sdm_demo_todolist/shared/resp"
	"sdm_demo_todolist/sqlx/dbal"
	"sdm_demo_todolist/sqlx/dbal/dto"

	"github.com/gin-gonic/gin"
)

type projectHandlers struct {
	dao *dbal.ProjectsDao
}

func NewProjectHandlers() shared.ProjectHandlers {
	return &projectHandlers{
		dao: dbal.NewProjectsDao(),
	}
}

// ProjectCreate
//
//	@Summary	create project
//	@Tags		Projects
//	@Id			ProjectCreate
//	@Accept		json
//	@Success	201
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/projects [post]
//	@Param		json	body	request.Project	true	"project data"
func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := h.dao.CreateProject(ctx, &dto.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// ProjectsReadAll
//
//	@Summary	get project list
//	@Tags		Projects
//	@Id			ProjectsReadAll
//	@Produce	json
//	@Success	200	{object}	[]dto.ProjectLi	"project list"
//	@Failure	500
//	@Security	none
//	@Router		/projects [get]
func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	projects, err := h.dao.GetProjects(ctx)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, projects)
}

// ProjectRead
//
//	@Summary	get project
//	@Tags		Projects
//	@Id			ProjectRead
//	@Produce	json
//	@Success	200	{object}	dto.Project	"project data"
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Security	none
//	@Router		/projects/{p_id} [get]
//	@Param		p_id	path	integer	true	"project id"
func (h *projectHandlers) ProjectRead(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	pr, err := h.dao.ReadProject(ctx, uri.PId)
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

// ProjectUpdate
//
//	@Summary	update project
//	@Tags		Projects
//	@Id			ProjectUpdate
//	@Accept		json
//	@Success	200
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/projects/{p_id} [put]
//	@Param		p_id	path	integer			true	"project id"
//	@Param		json	body	request.Project	true	"project data"
func (h *projectHandlers) ProjectUpdate(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if _, err := h.dao.UpdateProject(ctx, &dto.Project{PID: uri.PId, PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
	}
}

// ProjectDelete
//
//	@Summary	delete project
//	@Tags		Projects
//	@Id			ProjectDelete
//	@Success	204
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/projects/{p_id} [delete]
//	@Param		p_id	path	integer	true	"project id"
func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := h.dao.DeleteProject(ctx, &dto.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
