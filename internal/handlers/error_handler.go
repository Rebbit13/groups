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

func (e *ErrorHandler) parseContextToZapFields(context *gin.Context) []zap.Field {
	return []zap.Field{
		zap.String("ip", context.ClientIP()),
		zap.String("url", context.FullPath()),
	}
}

func (e *ErrorHandler) InternalServerError(context *gin.Context, message string) {
	e.logger.Error(
		fmt.Sprintf("500 INTERNAL SERVER ERROR response message : %s", message),
		e.parseContextToZapFields(context)...,
	)
	context.Status(500)
}

func (e *ErrorHandler) NotFound(context *gin.Context, message string) {
	e.logger.Info("404 NotFound response", e.parseContextToZapFields(context)...)
	context.JSON(404, e.jsonMessage(message))
}

func (e *ErrorHandler) BadRequest(context *gin.Context, message string) {
	e.logger.Info("400 BadRequest response", e.parseContextToZapFields(context)...)
	context.JSON(400, e.jsonMessage(message))
}

func (e *ErrorHandler) HandleError(context *gin.Context, err error) {
	switch err.(type) {
	case *service.GroupToAttachNotExistError:
		e.BadRequest(context, err.Error())
	case *service.RecursiveGroupDependenciesError:
		e.BadRequest(context, err.Error())
	case *service.RecordNotFound:
		e.NotFound(context, err.Error())
	default:
		e.InternalServerError(context, err.Error())
	}
}
