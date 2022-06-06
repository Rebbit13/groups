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

func BindGroupHandler(service GroupService, logger *zap.Logger, router *gin.Engine) {
	handler := &Group{service: service, logger: logger, errorHandler: &ErrorHandler{logger: logger}}
	router.GET("/group/:id", handler.Get)
	router.GET("/group", handler.GetAll)
	router.POST("/group", handler.Create)
	router.PATCH("/group/:id", handler.Update)
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

func (g *Group) Get(context *gin.Context) {
	id, err := g.parseIdFromContext(context)
	if err != nil {
		return
	}
	group, err := g.service.Get(id)
	if err != nil {
		g.errorHandler.HandleError(context, err)
		return
	}
	context.JSON(200, group)
}

func (g *Group) GetAll(context *gin.Context) {
	groups, err := g.service.GetAll()
	if err != nil {
		g.errorHandler.HandleError(context, err)
		return
	}
	context.JSON(200, groups)
}

func (g *Group) Update(context *gin.Context) {
	id, err := g.parseIdFromContext(context)
	if err != nil {
		return
	}
	group := g.bindContestToGroupToUpdateAndCreate(context)
	if group == nil {
		return
	}
	groupToUpdate := group.convertToGORMModel()
	groupToUpdate.ID = id
	groupUpdated, err := g.service.Update(groupToUpdate)
	if err != nil {
		g.errorHandler.HandleError(context, err)
		return
	}
	context.JSON(200, groupUpdated)
}

func (g *Group) Delete(context *gin.Context) {
	id, err := g.parseIdFromContext(context)
	if err != nil {
		return
	}
	err = g.service.Delete(id)
	if err != nil {
		g.errorHandler.HandleError(context, err)
		return
	}
	context.Status(204)
}

func (g *Group) Members(context *gin.Context) {
	id, err := g.parseIdFromContext(context)
	if err != nil {
		return
	}
	flat, err := strconv.ParseBool(context.DefaultQuery("flat", "true"))
	if err != nil {
		g.errorHandler.BadRequest(context, "flat param must be bool")
	}
	members, err := g.service.Members(id, flat)
	if err != nil {
		g.errorHandler.HandleError(context, err)
	}
	context.JSON(200, members)
}
