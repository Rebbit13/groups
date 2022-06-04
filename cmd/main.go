package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"groups/internal/handlers"
	"groups/internal/models"
	"groups/internal/service"
	"groups/internal/storage"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	entities := []interface{}{&models.Group{}, &models.Human{}}
	db := storage.NewSqliteDatabase(entities)

	groupService := service.NewGroup(db, logger)
	humanService := service.NewHuman(db, logger)

	router := gin.Default()

	handlers.BindGroupHandler(groupService, logger, router)
	handlers.BindHumanHandler(humanService, logger, router)

	if err = router.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}
