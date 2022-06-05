package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string `json:"name" binding:"required"`
	SuperGroupID *uint
	SubGroups    []*Group `gorm:"foreignkey:SuperGroupID"`
}
