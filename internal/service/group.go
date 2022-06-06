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
	ids := make([]uint, 0)
	for _, subGroup := range group.SubGroups {
		ids = append(ids, subGroup.ID)
	}
	err := g.checkRecursionDependence(group, ids)
	if err != nil {
		return err
	}
	foundedGroups := make([]*models.Group, 0)
	g.db.Where(ids).Find(&foundedGroups)
	if len(ids) != len(foundedGroups) {
		return NewGroupToAttachNotExistError("groups to attach must exist")
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
			return NewRecursiveGroupDependenciesError("can not add because of recursive group dependence")
		}
		supergroupId = superGroup.SuperGroupID
	}
	return nil
}

func (g *Group) getAllSubgroups(group *models.Group) ([]*models.Group, error) {
	var groups []*models.Group
	var subGroups []*models.Group
	result := g.db.Where("super_group_id = ?", group.ID).Find(&subGroups)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	for _, subgroup := range subGroups {
		groupsToAdd, err := g.getAllSubgroups(subgroup)
		if err != nil {
			return nil, err
		}
		groups = append(groups, groupsToAdd...)
	}
	return append(groups, group), nil
}

func (g *Group) Get(id uint) (*models.Group, error) {
	var group *models.Group
	result := g.db.Where("id = ?", id).First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, NewRecordNotFound(fmt.Sprintf("There is no group with id %d", id))
		}
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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return NewRecordNotFound(fmt.Sprintf("There is no group with id %d", id))
	}
	return result.Error
}

func (g *Group) Members(id uint, flat bool) ([]*models.Human, error) {
	ides := make([]uint, 0)
	if flat {
		group, err := g.Get(id)
		if err != nil {
			return nil, err
		}
		ides = []uint{group.ID}
	} else {
		group, err := g.Get(id)
		if err != nil {
			return nil, err
		}
		groups, err := g.getAllSubgroups(group)
		if err != nil {
			return nil, err
		}
		for _, g := range groups {
			ides = append(ides, g.ID)
		}
	}
	members := make([]*models.Human, 0)
	result := g.db.Table("humans").Where("id IN (SELECT human_id FROM humans_groups WHERE group_id IN ?)", ides).Scan(&members)
	if result.Error != nil {
		g.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return members, nil
}
