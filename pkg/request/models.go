package request

import (
	"sdm_demo_todolist/pkg/resp"

	"github.com/gin-gonic/gin"
)

type ProjectUri struct {
	PId int64 `uri:"p_id"  binding:"required,gte=1,lte=512"`
}

func BindProjectUri(ctx *gin.Context) (*ProjectUri, error) {
	var uri ProjectUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		resp.Abort400BadUri(ctx, err)
		return nil, err
	}
	return &uri, nil
}

type Project struct {
	PName string `json:"p_name" binding:"required,lte=256"`
}

type TaskUri struct {
	TId int64 `uri:"t_id" binding:"required"`
}

func BindTaskUri(ctx *gin.Context) (*TaskUri, error) {
	var uri TaskUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		resp.Abort400BadUri(ctx, err)
		return nil, err
	}
	return &uri, nil
}

type NewTask struct {
	TSubject string `json:"t_subject" binding:"required,gte=1,lte=512"`
}
