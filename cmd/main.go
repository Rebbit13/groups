package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"groups/internal/handlers"
	"groups/internal/service"
	"groups/internal/storage"
	"log"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	groupService := service.NewGroup(db, logger)
	humanService := service.NewHuman(db, logger)

	router := gin.Default()

	handlers.BindGroupHandler(groupService, logger, router)
	handlers.BindHumanHandler(humanService, logger, router)

	if err = router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT"))); err != nil {
		log.Fatal(err)
	}
}
