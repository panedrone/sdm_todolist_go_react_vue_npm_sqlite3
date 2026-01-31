package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"sdm_demo_todolist/internal/sqlx/dbal"
	"sdm_demo_todolist/internal/sqlx/dbal/dto"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
)

type projectHandlers struct {
}

func NewSqlxProjectHandlers() etc.ProjectHandlers {
	return &projectHandlers{}
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
	if err := dbal.CreateProject(ctx, &dto.Project{PName: req.PName}); err != nil {
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
	projects, err := dbal.GetProjects(ctx)
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
	if _, err := dbal.UpdateProject(ctx, &dto.Project{PID: uri.PId, PName: req.PName}); err != nil {
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
	if _, err := dbal.DeleteProject(ctx, &dto.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
