package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"groups/internal/models"
	"groups/pkg/utils"
	"log"
	"time"
)

func GenerateTestData(db *gorm.DB) {
	groups := make([]*models.Group, 0)
	humans := make([]*models.Human, 0)
	for i := 0; i < 10; i++ {
		group := &models.Group{Name: fmt.Sprintf("group_%s", utils.RandName(5))}
		result := db.Create(group)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		groups = append(groups, group)
	}
	for _, group := range groups {
		for i := 0; i < 3; i++ {
			human := &models.Human{
				Name:      fmt.Sprintf("name_%s", utils.RandName(5)),
				Surname:   fmt.Sprintf("surname_%s", utils.RandName(5)),
				Birthdate: time.Date(2005, 12, 12, 0, 0, 0, 0, time.UTC),
				Groups:    []*models.Group{group},
			}
			result := db.Create(human)
			if result.Error != nil {
				log.Fatal(result.Error)
			}
			humans = append(humans, human)
		}
	}
}
