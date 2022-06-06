package handlers

import (
	"gorm.io/gorm"
	"groups/internal/models"
	"time"
)

type GroupToUpdateAndCreate struct {
	Name   string `json:"name" binding:"required"`
	Groups []uint `json:"sub_groups" binding:"required"`
}

func (g *GroupToUpdateAndCreate) convertToGORMModel() *models.Group {
	groups := make([]*models.Group, 0)
	for _, id := range g.Groups {
		groups = append(groups, &models.Group{Model: gorm.Model{ID: id}})
	}
	return &models.Group{Name: g.Name, SubGroups: groups}
}

type HumanToUpdateAndCreate struct {
	Name      string    `json:"name" binding:"required"`
	Surname   string    `json:"surname" binding:"required"`
	Birthdate time.Time `json:"birthdate" binding:"required"`
	Groups    []uint    `json:"groups" binding:"required"`
}

func (h *HumanToUpdateAndCreate) convertToGORMModel() *models.Human {
	groups := make([]*models.Group, 0)
	for _, id := range h.Groups {
		groups = append(groups, &models.Group{Model: gorm.Model{ID: id}})
	}
	return &models.Human{
		Name:      h.Name,
		Surname:   h.Surname,
		Birthdate: h.Birthdate,
		Groups:    groups,
	}
}
