package main

import (
	"__elastic/pkg/common/container"
	"__elastic/pkg/common/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	c := container.New()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/db", func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("term")

		var result []models.Movie

		qb := c.DB.Model(&models.Movie{})

		if term != "" {
			qb = qb.Where("(title like \"%?%\") or (overview like \"%?%\") or (director like \"%?%\")", term, term, term)
		}

		start := time.Now()
		_ = qb.Find(&result)
		log.Printf("Took the database: %v", time.Since(start))

		marshalData, _ := json.Marshal(result)

		w.Write(marshalData)
	})

	r.Get("/es", func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("term")

		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"title": map[string]interface{}{
						"query":     term,
						"fuzziness": 2,
					},
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}
		res, err := c.ES.Search(
			c.ES.Search.WithContext(context.Background()),
			c.ES.Search.WithIndex("movies"),
			c.ES.Search.WithBody(&buf),
			c.ES.Search.WithTrackTotalHits(true),
			c.ES.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		var rBuf map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&rBuf); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		marshalData, _ := json.Marshal(rBuf)

		w.Write(marshalData)
	})

	log.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
