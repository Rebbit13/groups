package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"groups/internal/service"
)

type ErrorHandler struct {
	logger *zap.Logger
}

// json_message returns map[string]string{"message": value}
func (e *ErrorHandler) jsonMessage(value string) map[string]string {
	return map[string]string{"message": value}
}

func (e *ErrorHandler) parseContextToZapFields(ctx *gin.Context) []zap.Field {
	return []zap.Field{
		zap.String("ip", ctx.ClientIP()),
		zap.String("url", ctx.FullPath()),
	}
}

func (e *ErrorHandler) InternalServerError(ctx *gin.Context, message string) {
	e.logger.Error(
		fmt.Sprintf("500 INTERNAL SERVER ERROR response message : %s", message),
		e.parseContextToZapFields(ctx)...,
	)
	ctx.Status(500)
}

func (e *ErrorHandler) NotFound(ctx *gin.Context, message string) {
	e.logger.Info("404 NotFound response", e.parseContextToZapFields(ctx)...)
	ctx.JSON(404, e.jsonMessage(message))
}

func (e *ErrorHandler) BadRequest(ctx *gin.Context, message string) {
	e.logger.Info("400 BadRequest response", e.parseContextToZapFields(ctx)...)
	ctx.JSON(400, e.jsonMessage(message))
}

func (e *ErrorHandler) HandleError(ctx *gin.Context, err error) {
	switch err.(type) {
	case *service.GroupToAttachNotExistError:
		e.BadRequest(ctx, err.Error())
	case *service.RecursiveGroupDependenciesError:
		e.BadRequest(ctx, err.Error())
	case *service.RecordNotFound:
		e.NotFound(ctx, err.Error())
	default:
		e.InternalServerError(ctx, err.Error())
	}
}
