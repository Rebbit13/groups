package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func NewSqliteDatabase(entities []interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	for _, model := range entities {
		err = db.AutoMigrate(model)
		if err != nil {
			log.Fatal("failed to migrate schema")
		}
	}
	return db
}
