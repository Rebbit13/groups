package handlers

import "groups/internal/models"

type GroupToUpdateAndCreate struct {
	Name string
}

func (g *GroupToUpdateAndCreate) convertToGORMModel() *models.Group {
	return &models.Group{Name: g.Name}
}
