package main

import (
	"__elastic/pkg/common/container"
	"__elastic/pkg/common/models"
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
		log.Printf("Took the database: %v", time.Now().Sub(start))

		marshalData, _ := json.Marshal(result)

		w.Write(marshalData)
	})

	http.ListenAndServe(":3000", r)
}
