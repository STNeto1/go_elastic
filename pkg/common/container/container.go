package container

import (
	"__elastic/pkg/common/db"

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
