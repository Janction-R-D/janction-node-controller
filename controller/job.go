package controller

import (
	"janction/logic"
	"janction/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetJobHandler(c *gin.Context) {
	var params model.FormGetJob
	if err := c.ShouldBindQuery(&params); err != nil {
		zap.L().Error("Invalid query params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	job, err := logic.GetJob(params.NodeID)
	if err != nil {
		zap.L().Error("GetJob failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, job)
}
