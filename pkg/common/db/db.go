package db

import (
	"__elastic/pkg/common/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./assets/movies.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error creating", err)
	}

	err = db.AutoMigrate(&models.Movie{})
	if err != nil {
		log.Fatalln("Error migrating", err)
	}

	return db
}
