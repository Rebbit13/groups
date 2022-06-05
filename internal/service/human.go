package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"groups/internal/models"
)

type Human struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewHuman(db *gorm.DB, logger *zap.Logger) *Human {
	return &Human{db: db, logger: logger}
}

func (h *Human) checkGroupsToAttach(human *models.Human) error {
	if len(human.Groups) == 0 {
		return nil
	}
	ids := []uint{}
	foundedGroups := []*models.Group{}
	for _, group := range human.Groups {
		ids = append(ids, group.ID)
	}
	h.db.Where(ids).Find(&foundedGroups)
	if len(ids) != len(foundedGroups) {
		return errors.New("groups to attach must exist")
	}
	return nil
}

func (h *Human) Get(id uint) (*models.Human, error) {
	var human *models.Human
	result := h.db.Preload("Groups").Where("id = ?", id).First(&human)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return human, nil
}

func (h *Human) GetAll() ([]*models.Human, error) {
	var humans []*models.Human
	result := h.db.Preload("Groups").Find(&humans)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return humans, nil
}

func (h *Human) Create(human *models.Human) (*models.Human, error) {
	err := h.checkGroupsToAttach(human)
	if err != nil {
		return nil, err
	}
	result := h.db.Create(human)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return human, nil
}

func (h *Human) Update(human *models.Human) (*models.Human, error) {
	foundedHuman, err := h.Get(human.ID)
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	err = h.checkGroupsToAttach(human)
	if err != nil {
		return nil, err
	}
	result := h.db.Model(&foundedHuman).Updates(&human)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	err = h.db.Model(&foundedHuman).Association("Groups").Replace(human.Groups)
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	return foundedHuman, err
}

func (h *Human) Delete(id uint) error {
	var human *models.Human
	result := h.db.Where("id = ?", id).Delete(&human)
	return result.Error
}
