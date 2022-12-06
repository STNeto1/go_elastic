package db

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitES() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalln("Error starting ES", err)
	}

	return es
}
