package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	_ "groups/docs"
	"groups/internal/handlers"
	"groups/internal/service"
	"groups/internal/storage"
	"log"
	"os"
	"strconv"
)

// @title Groups Api
// @version 2.0
// @description This is a simple server for create groups and humans
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
// @schemes http
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

	testData := os.Getenv("TEST_DATA")
	if testData == "" {
		testData = "false"
	}
	ifNeedTestData, err := strconv.ParseBool(os.Getenv("TEST_DATA"))
	if err != nil {
		log.Fatal(err)
	}

	// generates test data in db
	if ifNeedTestData {
		handlers.GenerateTestData(db)
	}

	groupService := service.NewGroup(db, logger)
	humanService := service.NewHuman(db, logger)

	r := gin.Default()

	v1 := r.Group("/v1")
	handlers.BindGroupHandler(groupService, logger, v1)
	handlers.BindHumanHandler(humanService, logger, v1)

	r.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err = r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
		log.Fatal(err)
	}
}
