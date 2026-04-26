package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sdm_demo_todolist/internal/gorm/dbal"
	"sdm_demo_todolist/internal/gorm/dbal/models"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/datetime"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type taskHandlers struct {
}

func NewGormTaskHandlers() etc.TaskHandlers {
	return &taskHandlers{}
}

// TaskCreate godoc
// @Summary Create a new task
// @Description Create a new task for a specific project
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param p_id path int true "Project ID"
// @Param task body request.NewTask true "Task Subject"
// @Success 201
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects/{p_id}/tasks [post]
func (h *taskHandlers) TaskCreate(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	var req request.NewTask
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	t := models.Task{}
	t.PID = uri.PId
	t.TSubject = req.TSubject
	t.TPriority = 1
	t.TDate = datetime.NowLocalString()
	if err := dbal.CreateTask(ctx, &t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// TaskRead godoc
// @Summary Get a task by ID
// @Description Get details of a task by its ID
// @Tags tasks
// @Produce  json
// @Param t_id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /tasks/{t_id} [get]
func (h *taskHandlers) TaskRead(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	task, err := dbal.ReadTask(ctx, uri.TId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Abort404NotFound(ctx, err)
			return
		}
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, task)
}

// TasksReadByProject godoc
// @Summary Get tasks by project
// @Description Get all tasks associated with a specific project
// @Tags tasks
// @Produce  json
// @Param p_id path int true "Project ID"
// @Success 200 {array} models.Task
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /projects/{p_id}/tasks [get]
func (h *taskHandlers) TasksReadByProject(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	// tasks, err := tDao.RawProjectTasks(ctx, uri.PId)
	tasks, err := dbal.ReadProjectTasks(ctx, uri.PId)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, tasks)
}

// TaskUpdate godoc
// @Summary Update a task
// @Description Update an existing task's details by its ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param t_id path int true "Task ID"
// @Param task body models.Task true "Task Data"
// @Success 200
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /tasks/{t_id} [put]
func (h *taskHandlers) TaskUpdate(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	t, err := dbal.ReadTask(ctx, uri.TId)
	if err != nil {
		resp.Abort400hBadRequest(ctx, err.Error())
		return
	}
	var req models.Task
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
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Param t_id path int true "Task ID"
// @Success 204
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /tasks/{t_id} [delete]
func (h *taskHandlers) TaskDelete(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	if _, err := dbal.DeleteTask(ctx, &models.Task{TID: uri.TId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
