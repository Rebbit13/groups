package handlers

import "groups/internal/models"

type GroupService interface {
	Get(uint) (*models.Group, error)
	GetAll() ([]*models.Group, error)
	Create(*models.Group) (*models.Group, error)
	Update(*models.Group) (*models.Group, error)
	Delete(uint) error
	Members(uint, bool) ([]*models.Human, error)
}

type HumanService interface {
	Get(uint) (*models.Human, error)
	GetAll() ([]*models.Human, error)
	Create(*models.Human) (*models.Human, error)
	Update(*models.Human) (*models.Human, error)
	Delete(uint) error
}
