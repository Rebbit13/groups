package models

import (
	"gorm.io/gorm"
	"time"
)

type Human struct {
	gorm.Model
	Name      string    `json:"name" binding:"required"`
	Surname   string    `json:"surname" binding:"required"`
	Birthdate time.Time `json:"birthdate" binding:"required"`
	Groups    []*Group  `gorm:"many2many:humans_groups;" json:"groups" binding:"required"`
}

func (h *Human) TableName() string {
	return "humans"
}
