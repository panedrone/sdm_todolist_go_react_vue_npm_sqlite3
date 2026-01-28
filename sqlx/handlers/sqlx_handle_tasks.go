package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sdm_demo_todolist/pkg"
	"sdm_demo_todolist/pkg/datetime"
	"sdm_demo_todolist/pkg/request"
	"sdm_demo_todolist/pkg/resp"
	"sdm_demo_todolist/sqlx/dbal"
	"sdm_demo_todolist/sqlx/dbal/dto"

	"github.com/gin-gonic/gin"
)

type taskHandlers struct {
	dao *dbal.TasksDao
}

func NewSqlxTaskHandlers() etc.TaskHandlers {
	return &taskHandlers{
		dao: dbal.NewTasksDao(),
	}
}

// TaskCreate
//
//	@Summary	create task
//	@Tags		Tasks
//	@Id			TaskCreate
//	@Accept		json
//	@Param		json	body	request.NewTask	true	"task data"
//	@Success	201
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/projects/{p_id}/tasks [post]
//	@Param		p_id	path	integer	true	"project id"
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
	if err := h.dao.CreateTask(ctx, &t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// TasksReadByProject
//
//	@Summary	get project tasks
//	@Tags		Tasks
//	@Id			TasksReadByProject
//	@Produce	json
//	@Success	200	{object}	[]dto.TaskLi	"project tasks"
//	@Failure	500
//	@Security	none
//	@Router		/projects/{p_id}/tasks [get]
//	@Param		p_id	path	integer	true	"project id"
func (h *taskHandlers) TasksReadByProject(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	tasks, err := h.dao.GetGroupTasks(ctx, uri.PId)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, tasks)
}

// TaskRead
//
//	@Summary	get task
//	@Tags		Tasks
//	@Id			TaskRead
//	@Produce	json
//	@Success	200	{object}	dto.Task	"task data"
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Security	none
//	@Router		/tasks/{t_id} [get]
//	@Param		t_id	path	integer	true	"task id"
func (h *taskHandlers) TaskRead(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	task, err := h.dao.ReadTask(ctx, uri.TId)
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

// TaskUpdate
//
//	@Summary	update task
//	@Tags		Tasks
//	@Id			TaskUpdate
//	@Accept		json
//	@Success	200
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/tasks/{t_id} [put]
//	@Param		t_id	path	integer		true	"task id"
//	@Param		json	body	dto.Task	true	"task data"
func (h *taskHandlers) TaskUpdate(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	t, err := h.dao.ReadTask(ctx, uri.TId)
	if err != nil {
		resp.Abort400BadUri(ctx, err)
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
	if _, err = h.dao.UpdateTask(ctx, t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
}

// TaskDelete
//
//	@Summary	delete task
//	@Tags		Tasks
//	@Id			TaskDelete
//	@Success	204
//	@Failure	400
//	@Failure	500
//	@Security	none
//	@Router		/tasks/{t_id} [delete]
//	@Param		t_id	path	integer	true	"task id"
func (h *taskHandlers) TaskDelete(ctx *gin.Context) {
	uri, err := request.BindTaskUri(ctx)
	if err != nil {
		return
	}
	t := dto.Task{TID: uri.TId}
	if _, err := h.dao.DeleteTask(ctx, &t); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
