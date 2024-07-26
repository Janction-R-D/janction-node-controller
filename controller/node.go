package controller

import (
	"janction/logic"
	"janction/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterNodeHandler(c *gin.Context) {
	var params model.FormRegisterNode
	if err := c.ShouldBindJSON(&params); err != nil {
		zap.L().Error("Invalid request body", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := logic.RegisterNode(&params); err != nil {
		zap.L().Error("RegisterNode failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, http.StatusOK)
}

func PingHandler(c *gin.Context) {
	var params model.FormPing
	if err := c.ShouldBindJSON(&params); err != nil {
		zap.L().Error("Invalid request body", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := logic.Ping(params.NodeID); err != nil {
		zap.L().Error("Ping failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, http.StatusOK)
}
