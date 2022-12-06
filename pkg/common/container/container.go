package container

import (
	"__elastic/pkg/common/db"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB
	ES *elasticsearch.Client
}

func New() *Container {
	dbConn := db.Init()
	esClient := db.InitES()

	c := Container{
		DB: dbConn,
		ES: esClient,
	}

	return &c
}
