package service

import (
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

func (h *Human) Get(id uint) (*models.Human, error) {
	var human *models.Human
	result := h.db.Where("id = ?", id).First(&human)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return human, nil
}

func (h *Human) GetAll() ([]*models.Human, error) {
	var humans []*models.Human
	result := h.db.Find(&humans)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return humans, nil
}

func (h *Human) Create(human *models.Human) (*models.Human, error) {
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
	result := h.db.Model(&foundedHuman).Updates(&human)
	if result.Error != nil {
		h.logger.Error(result.Error.Error())
		return nil, result.Error
	}
	return foundedHuman, err
}

func (h *Human) Delete(id uint) error {
	var human *models.Human
	result := h.db.Where("id = ?", id).Delete(&human)
	return result.Error
}
