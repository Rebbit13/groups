package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"groups/internal/models"
)

type Group struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGroup(db *gorm.DB, logger *zap.Logger) *Group {
	return &Group{db: db, logger: logger}
}

func (g *Group) Get(id uint) (*models.Group, error) {
	var group *models.Group
	result := g.db.Where("id = ?", id).First(&group)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return group, nil
}

func (g *Group) GetAll() ([]*models.Group, error) {
	var groups []*models.Group
	result := g.db.Find(&groups)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return groups, nil
}

func (g *Group) Create(group *models.Group) error {
	result := g.db.Create(group)
	return result.Error
}

func (g *Group) Update(group *models.Group) (*models.Group, error) {
	foundedGroup, err := g.Get(group.ID)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}
	result := g.db.Model(&foundedGroup).Updates(&group)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return foundedGroup, err
}

func (g *Group) Delete(id uint) error {
	var group *models.Group
	result := g.db.Where("id = ?", id).Delete(&group)
	return result.Error
}
