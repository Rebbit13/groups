package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Human struct {
	service      HumanService
	logger       *zap.Logger
	errorHandler *ErrorHandler
}

func BindHumanHandler(service HumanService, logger *zap.Logger, router *gin.RouterGroup) {
	handler := &Human{service: service, logger: logger, errorHandler: &ErrorHandler{logger: logger}}
	router.GET("/human/:id", handler.Get)
	router.GET("/human", handler.GetAll)
	router.POST("/human", handler.Create)
	router.PATCH("/human/:id", handler.Update)
	router.DELETE("/human/:id", handler.Delete)
}

func (h *Human) parseIdFromContext(ctx *gin.Context) (uint, error) {
	humanId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler.BadRequest(ctx, "Can not parse id. Must be unsigned integer.")
		return 0, err
	}
	return uint(humanId), nil
}

func (h *Human) bindContestToHumanToUpdateAndCreate(ctx *gin.Context) *HumanToUpdateAndCreate {
	human := &HumanToUpdateAndCreate{}
	err := ctx.BindJSON(human)
	if err != nil {
		h.errorHandler.BadRequest(ctx, err.Error())
		return nil
	}
	return human
}

// Create godoc
// @Summary create human
// @Tags human
// @Accept json
// @Produce json
// @Param message  body  HumanToUpdateAndCreate  true  "Human"
// @Success 201 {object} map[string]interface{}
// @Failure 400  {object}  map[string]interface{}
// @Failure 500
// @Router /human [post]
func (h *Human) Create(ctx *gin.Context) {
	human := h.bindContestToHumanToUpdateAndCreate(ctx)
	if human == nil {
		return
	}
	humanCreated, err := h.service.Create(human.convertToGORMModel())
	if err != nil {
		h.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(201, humanCreated)
}

// Get godoc
// @Summary get one human by id
// @Tags human
// @Accept */*
// @Produce json
// @Param human_id path int true "Human ID"
// @Success 200 {array} map[string]string
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /human/{human_id} [get]
func (h *Human) Get(ctx *gin.Context) {
	id, err := h.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	human, err := h.service.Get(id)
	if err != nil {
		h.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, human)
}

// GetAll godoc
// @Summary get all humans
// @Tags human
// @Accept */*
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500
// @Router /human [get]
func (h *Human) GetAll(ctx *gin.Context) {
	humans, err := h.service.GetAll()
	if err != nil {
		h.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, humans)
}

// Update godoc
// @Summary update one human by id
// @Tags human
// @Accept */*
// @Produce json
// @Param human_id path int true "Human ID"
// @Param message  body  HumanToUpdateAndCreate  true  "Human"
// @Success 200 {array} map[string]string
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /human/{human_id} [put]
func (h *Human) Update(ctx *gin.Context) {
	id, err := h.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	human := h.bindContestToHumanToUpdateAndCreate(ctx)
	if human == nil {
		return
	}
	humanToUpdate := human.convertToGORMModel()
	humanToUpdate.ID = id
	humanUpdated, err := h.service.Update(humanToUpdate)
	if err != nil {
		h.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, humanUpdated)
}

// Delete godoc
// @Summary delete one human by id
// @Tags human
// @Accept */*
// @Produce */*
// @Param human_id path int true "Human ID"
// @Success 204
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /human/{human_id} [delete]
func (h *Human) Delete(ctx *gin.Context) {
	id, err := h.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		h.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.Status(204)
}
