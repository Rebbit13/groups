package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Human struct {
	service HumanService
	logger  *zap.Logger
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
		context.Status(500)
		return
	}
	context.JSON(201, humanCreated)
}

func (h *Human) Get(context *gin.Context) {
	humanId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	human, err := h.service.Get(uint(humanId))
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "404"))
		context.JSON(404, json_message(fmt.Sprintf("There is no human with id %d", humanId)))
		return
	}
	context.JSON(200, human)
}

func (h *Human) GetAll(context *gin.Context) {
	humans, err := h.service.GetAll()
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(200, humans)
}

func (h *Human) Update(context *gin.Context) {
	humanId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	human := &HumanToUpdateAndCreate{}
	err = context.BindJSON(human)
	if err != nil {
		return
	}
	humanToUpdate := human.convertToGORMModel()
	humanToUpdate.ID = uint(humanId)
	humanUpdated, err := h.service.Update(humanToUpdate)
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(200, humanUpdated)
}

func (h *Human) Delete(context *gin.Context) {
	humanId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	err = h.service.Delete(uint(humanId))
	if err != nil {
		h.logger.Error(err.Error(), zap.String("status_code", "404"))
		context.JSON(404, json_message(fmt.Sprintf("There is no human with id %d", humanId)))
		return
	}
	context.Status(204)
}
