package service

import (
	"errors"
	"fmt"
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

func (g *Group) checkGroupsToAttach(group *models.Group) error {
	if len(group.SubGroups) == 0 {
		return nil
	}
	ids := []uint{}
	for _, subGroup := range group.SubGroups {
		ids = append(ids, subGroup.ID)
	}
	err := g.checkRecursionDependence(group, ids)
	if err != nil {
		return err
	}
	foundedGroups := []*models.Group{}
	g.db.Where(ids).Find(&foundedGroups)
	if len(ids) != len(foundedGroups) {
		return errors.New("groups to attach must exist")
	}
	return nil
}

func (g *Group) checkIfIDinSlice(id uint, slice []uint) bool {
	for _, candidate := range slice {
		if id == candidate {
			return true
		}
	}
	return false
}

func (g *Group) checkRecursionDependence(group *models.Group, subGroupsIdes []uint) error {
	for supergroupId := &group.ID; supergroupId != nil; {
		var superGroup *models.Group
		g.db.Where("id = ?", supergroupId).First(&superGroup)
		fmt.Println(superGroup)
		if g.checkIfIDinSlice(*supergroupId, subGroupsIdes) {
			return errors.New("can not add because of the recursive group dependence")
		}
		supergroupId = superGroup.SuperGroupID
	}
	return nil
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

func (g *Group) Create(group *models.Group) (*models.Group, error) {
	result := g.db.Create(group)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return group, nil
}

func (g *Group) Update(group *models.Group) (*models.Group, error) {
	foundedGroup, err := g.Get(group.ID)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}
	err = g.checkGroupsToAttach(group)
	if err != nil {
		return nil, err
	}
	result := g.db.Model(&foundedGroup).Updates(&group)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}

	err = g.db.Model(&foundedGroup).Association("SubGroups").Replace(group.SubGroups)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}

	return foundedGroup, err
}

func (g *Group) Delete(id uint) error {
	var group *models.Group
	result := g.db.Where("id = ?", id).Delete(&group)
	return result.Error
}
