package container

import (
	"__elastic/pkg/common/db"
	. "github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB
}

func New() *Container {
	db := db.Init()

	c := Container{
		DB: db,
	}

	return &c
}
