package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type Group struct {
	service GroupService
	logger  *zap.Logger
}

func (g *Group) Create(context *gin.Context) {
	group := &GroupToUpdateAndCreate{}
	err := context.BindJSON(group)
	if err != nil {
		return
	}
	groupCreated, err := g.service.Create(group.convertToGORMModel())
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(201, groupCreated)
}

func (g *Group) Get(context *gin.Context) {
	groupId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	group, err := g.service.Get(uint(groupId))
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "404"))
		context.JSON(404, json_message(fmt.Sprintf("There is no group with id %d", groupId)))
		return
	}
	context.JSON(200, group)
}

func (g *Group) GetAll(context *gin.Context) {
	groups, err := g.service.GetAll()
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(200, groups)
}

func (g *Group) Update(context *gin.Context) {
	groupId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	group := &GroupToUpdateAndCreate{}
	err = context.BindJSON(group)
	if err != nil {
		return
	}
	groupToUpdate := group.convertToGORMModel()
	groupToUpdate.ID = uint(groupId)
	groupUpdated, err := g.service.Update(groupToUpdate)
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "500"))
		context.Status(500)
		return
	}
	context.JSON(200, groupUpdated)
}

func (g *Group) Delete(context *gin.Context) {
	groupId, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "400"))
		context.JSON(400, json_message("Can not parse id. Must be unsigned integer."))
		return
	}
	err = g.service.Delete(uint(groupId))
	if err != nil {
		g.logger.Error(err.Error(), zap.String("status_code", "404"))
		context.JSON(404, json_message(fmt.Sprintf("There is no group with id %d", groupId)))
		return
	}
	context.Status(204)
}
