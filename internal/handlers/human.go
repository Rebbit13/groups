package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Human struct {
	service      HumanService
	logger       *zap.Logger
	errorHandler *ErrorHandler
}

func BindHumanHandler(service HumanService, logger *zap.Logger, router *gin.Engine) {
	handler := &Human{service: service, logger: logger, errorHandler: &ErrorHandler{logger: logger}}
	router.GET("/human/:id", handler.Get)
	router.GET("/human", handler.GetAll)
	router.POST("/human", handler.Create)
	router.PATCH("/human/:id", handler.Update)
	router.DELETE("/human/:id", handler.Delete)
}

func (h *Human) parseIdFromContext(context *gin.Context) (uint, error) {
	humanId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorHandler.BadRequest(context, "Can not parse id. Must be unsigned integer.")
		return 0, err
	}
	return uint(humanId), nil
}

func (h *Human) Create(context *gin.Context) {
	human := &HumanToUpdateAndCreate{}
	err := context.BindJSON(human)
	if err != nil {
		return
	}
	humanCreated, err := h.service.Create(human.convertToGORMModel())
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "500"))
		h.errorHandler.InternalServerError(context)
		return
	}
	context.JSON(201, humanCreated)
}

func (h *Human) Get(context *gin.Context) {
	id, err := h.parseIdFromContext(context)
	if err != nil {
		return
	}
	human, err := h.service.Get(id)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorHandler.NotFound(context, fmt.Sprintf("There is no human with id %d", id))
		return
	}
	context.JSON(200, human)
}

func (h *Human) GetAll(context *gin.Context) {
	humans, err := h.service.GetAll()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorHandler.InternalServerError(context)
		return
	}
	context.JSON(200, humans)
}

func (h *Human) Update(context *gin.Context) {
	id, err := h.parseIdFromContext(context)
	if err != nil {
		return
	}
	human := &HumanToUpdateAndCreate{}
	err = context.BindJSON(human)
	if err != nil {
		return
	}
	humanToUpdate := human.convertToGORMModel()
	humanToUpdate.ID = id
	humanUpdated, err := h.service.Update(humanToUpdate)
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(200, humanUpdated)
}

func (h *Human) Delete(context *gin.Context) {
	id, err := h.parseIdFromContext(context)
	if err != nil {
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorHandler.NotFound(context, fmt.Sprintf("There is no human with id %d", id))
		return
	}
	context.Status(204)
}
