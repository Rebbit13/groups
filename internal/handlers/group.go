package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Group struct {
	service      GroupService
	logger       *zap.Logger
	errorHandler *ErrorHandler
}

func NewGroupHandler(service GroupService, logger *zap.Logger) *Group {
	return &Group{service: service, logger: logger, errorHandler: &ErrorHandler{logger: logger}}
}

func BindGroupHandler(service GroupService, logger *zap.Logger, router *gin.RouterGroup) {
	handler := &Group{service: service, logger: logger, errorHandler: &ErrorHandler{logger: logger}}
	router.GET("/group/:id", handler.Get)
	router.GET("/group", handler.GetAll)
	router.POST("/group", handler.Create)
	router.PUT("/group/:id", handler.Update)
	router.DELETE("/group/:id", handler.Delete)
	router.GET("/group/:id/members", handler.Members)
}

func (g *Group) parseIdFromContext(context *gin.Context) (uint, error) {
	groupId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		g.logger.Error(err.Error())
		g.errorHandler.BadRequest(context, "Can not parse id. Must be unsigned integer.")
		return 0, err
	}
	return uint(groupId), nil
}

func (g *Group) bindContestToGroupToUpdateAndCreate(context *gin.Context) *GroupToUpdateAndCreate {
	group := &GroupToUpdateAndCreate{}
	err := context.BindJSON(group)
	if err != nil {
		g.errorHandler.BadRequest(context, err.Error())
		return nil
	}
	return group
}

// Create godoc
// @Summary create group
// @Tags group
// @Accept json
// @Produce json
// @Param message  body  GroupToUpdateAndCreate  true  "Group"
// @Success 201 {object} map[string]interface{}
// @Failure 400  {object}  map[string]interface{}
// @Failure 500
// @Router /group [post]
func (g *Group) Create(context *gin.Context) {
	group := g.bindContestToGroupToUpdateAndCreate(context)
	if group == nil {
		return
	}
	groupCreated, err := g.service.Create(group.convertToGORMModel())
	if err != nil {
		g.errorHandler.HandleError(context, err)
		return
	}
	context.JSON(201, groupCreated)
}

// Get godoc
// @Summary get one group by id
// @Tags group
// @Accept */*
// @Produce json
// @Param group_id path int true "Group ID"
// @Success 200 {array} map[string]string
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /group/{group_id} [get]
func (g *Group) Get(ctx *gin.Context) {
	id, err := g.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	group, err := g.service.Get(id)
	if err != nil {
		g.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, group)
}

// GetAll godoc
// @Summary get all groups
// @Tags group
// @Accept */*
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500
// @Router /group [get]
func (g *Group) GetAll(ctx *gin.Context) {
	groups, err := g.service.GetAll()
	if err != nil {
		g.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, groups)
}

// Update godoc
// @Summary update one group by id
// @Tags group
// @Accept */*
// @Produce json
// @Param group_id path int true "Group ID"
// @Param message  body  GroupToUpdateAndCreate  true  "Group"
// @Success 200 {array} map[string]string
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /group/{group_id} [put]
func (g *Group) Update(ctx *gin.Context) {
	id, err := g.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	group := g.bindContestToGroupToUpdateAndCreate(ctx)
	if group == nil {
		return
	}
	groupToUpdate := group.convertToGORMModel()
	groupToUpdate.ID = id
	groupUpdated, err := g.service.Update(groupToUpdate)
	if err != nil {
		g.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, groupUpdated)
}

// Delete godoc
// @Summary delete one group by id
// @Tags group
// @Accept */*
// @Produce */*
// @Param group_id path int true "Group ID"
// @Success 204
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /group/{group_id} [delete]
func (g *Group) Delete(ctx *gin.Context) {
	id, err := g.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	err = g.service.Delete(id)
	if err != nil {
		g.errorHandler.HandleError(ctx, err)
		return
	}
	ctx.Status(204)
}

// Members godoc
// @Summary get members of group by group id
// @Tags group
// @Accept */*
// @Produce json
// @Param group_id path int true "Group ID"
// @Param default query bool false "bool default" default(true)
// @Success 200 {array} map[string]interface{}
// @Failure 400  {object}  map[string]interface{}
// @Failure 404  {object}  map[string]interface{}
// @Failure 500
// @Router /group/{group_id}/members [get]
func (g *Group) Members(ctx *gin.Context) {
	id, err := g.parseIdFromContext(ctx)
	if err != nil {
		return
	}
	flat, err := strconv.ParseBool(ctx.DefaultQuery("flat", "true"))
	if err != nil {
		g.errorHandler.BadRequest(ctx, "flat param must be bool")
	}
	members, err := g.service.Members(id, flat)
	if err != nil {
		g.errorHandler.HandleError(ctx, err)
	}
	ctx.JSON(200, members)
}
