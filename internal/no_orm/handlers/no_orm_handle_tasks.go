package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sdm_demo_todolist/internal/no_orm/dbal"
	"sdm_demo_todolist/internal/no_orm/dbal/dto"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/datetime"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
)

type taskHandlers struct {
}

func NewNoOrmTaskHandlers() etc.TaskHandlers {
	return &taskHandlers{}
}

// TaskCreate godoc
// @Summary Create a new task (No-ORM)
// @Description Create a new task for a specific project using raw SQL
// @Tags tasks-no-orm
// @Accept  json
// @Produce  json
// @Param p_id path int true "Project ID"
// @Param task body request.NewTask true "Task Subject"
// @Success 201
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects/{p_id}/tasks [post]
func (h *taskHandlers) TaskCreate(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	var inTask request.NewTask
	if err := ctx.ShouldBindJSON(&inTask); err != nil {
		resp.Abort400hBadRequest(ctx, err.Error())
		return
	}
	t := dto.Task{}
	t.PID = uri.PId
	t.TSubject = inTask.TSubject
	t.TPriority = 1
	t.TDate = datetime.NowLocalString()
	if err := dbal.CreateTask(ctx, &t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// TasksReadByProject godoc
// @Summary Get tasks by project (No-ORM)
// @Description Get all tasks associated with a specific project using raw SQL
// @Tags tasks-no-orm
// @Produce  json
// @Param p_id path int true "Project ID"
// @Success 200 {array} dto.TaskLi
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/projects/{p_id}/tasks [get]
func (h *taskHandlers) TasksReadByProject(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	tasks, err := dbal.ReadByProject(ctx, uri.PId)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, tasks)
}

// TaskRead godoc
// @Summary Get a task by ID (No-ORM)
// @Description Get details of a task by its ID using raw SQL
// @Tags tasks-no-orm
// @Produce  json
// @Param t_id path int true "Task ID"
// @Success 200 {object} dto.Task
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/tasks/{t_id} [get]
func (h *taskHandlers) TaskRead(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	task, err := dbal.ReadTask(ctx, uri.TId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.Abort404NotFound(ctx, err)
			return
		}
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, task)
}

// TaskUpdate godoc
// @Summary Update a task (No-ORM)
// @Description Update an existing task's details by its ID using raw SQL
// @Tags tasks-no-orm
// @Accept  json
// @Produce  json
// @Param t_id path int true "Task ID"
// @Param task body dto.Task true "Task Data"
// @Success 200
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/tasks/{t_id} [put]
func (h *taskHandlers) TaskUpdate(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	t, err := dbal.ReadTask(ctx, uri.TId)
	if err != nil {
		resp.Abort404NotFound(ctx, err)
		return
	}
	var req dto.Task
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	errDetails := make(map[string]string)
	_, err = datetime.Parse(req.TDate)
	if err != nil {
		errDetails["t_date"] = fmt.Sprintf("date format expected '%s': %s", datetime.TimeFormat, err.Error())
	}
	if len(req.TSubject) == 0 {
		errDetails["t_subject"] = fmt.Sprintf("subject required")
	}
	if req.TPriority <= 0 {
		errDetails["t_priority"] = fmt.Sprintf("invalid priority: %v, an integer 1..10 expected", req.TPriority)
	}
	if len(errDetails) > 0 {
		resp.AbortWithError(ctx, &resp.ErrorResponse{
			Error:   "bad request",
			Details: errDetails,
		})
		return
	}
	t.TID = uri.TId
	t.PID = req.PID
	t.TSubject = req.TSubject
	t.TPriority = req.TPriority
	t.TDate = req.TDate
	t.TComments = req.TComments
	if _, err = dbal.UpdateTask(ctx, t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
}

// TaskDelete godoc
// @Summary Delete a task (No-ORM)
// @Description Delete a task by its ID using raw SQL
// @Tags tasks-no-orm
// @Param t_id path int true "Task ID"
// @Success 204
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /no_orm/tasks/{t_id} [delete]
func (h *taskHandlers) TaskDelete(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	t := dto.Task{TID: uri.TId}
	if _, err := dbal.DeleteTask(ctx, &t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
