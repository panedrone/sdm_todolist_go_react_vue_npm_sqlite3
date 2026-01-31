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
