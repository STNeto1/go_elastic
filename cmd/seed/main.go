package main

import (
	"__elastic/pkg/common/db"
	"__elastic/pkg/common/models"
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

type Shape struct {
	Poster       string  `json:"poster"`
	Title        string  `json:"title"`
	ReleaseYear  int32   `json:"release_year"`
	Certificate  string  `json:"certificate"`
	Runtime      string  `json:"runtime"`
	Genre        string  `json:"genre"`
	IMDBRating   float32 `json:"imdb_rating"`
	Overview     string  `json:"overview"`
	MetaScore    int32   `json:"meta_score"`
	Director     string  `json:"director"`
	Votes        int64   `json:"votes"`
	GrossRevenue int64   `json:"gross_revenue"`
}

func main() {
	dbConn := db.Init()
	es := db.InitES()

	f, err := os.Open("assets/dataset.csv")
	if err != nil {
		log.Fatalln("Error loading dataset", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalln("Error reading dataset", err)
	}

	movies := []Shape{}

	for i, line := range data {
		// header
		if i == 0 {
			continue
		}

		movie := Shape{
			Poster:       line[0],
			Title:        line[1],
			ReleaseYear:  parseInt32(line[2]),
			Certificate:  line[3],
			Runtime:      line[4],
			Genre:        line[5],
			IMDBRating:   parseFloat32(line[6]),
			Overview:     line[7],
			MetaScore:    parseInt32(line[8]),
			Director:     line[9],
			Votes:        parseInt64(line[14]),
			GrossRevenue: parseInt64(line[15]),
		}

		movies = append(movies, movie)
	}

	dbConn.Transaction(func(tx *gorm.DB) error {
		delRes := tx.Where("1 = 1").Delete(&models.Movie{})
		if delRes.Error != nil {
			return delRes.Error
		}

		for _, movie := range movies {
			res := tx.Create(&models.Movie{
				Poster:       movie.Poster,
				Title:        movie.Title,
				ReleaseYear:  movie.ReleaseYear,
				Certificate:  movie.Certificate,
				Runtime:      movie.Runtime,
				Genre:        movie.Genre,
				IMDBRating:   movie.IMDBRating,
				Overview:     movie.Overview,
				MetaScore:    movie.MetaScore,
				Director:     movie.Director,
				Votes:        movie.Votes,
				GrossRevenue: movie.GrossRevenue,
			})

			if res.Error != nil {
				return err
			}
		}

		return nil
	})

	var wg sync.WaitGroup
	for _, movie := range movies {
		wg.Add(1)

		go func(movie Shape) {
			data, err := json.Marshal(movie)
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}

			idx := strconv.Itoa(int(movie.ReleaseYear)) + movie.Title
			req := esapi.IndexRequest{
				Index:      "movies",
				DocumentID: idx,
				Body:       strings.NewReader(string(data)),
				Refresh:    "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%s", res.Status(), idx)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}

			wg.Done()
		}(movie)
	}

	wg.Wait()
}

func parseInt32(elem string) int32 {
	val, err := strconv.ParseInt(elem, 10, 32)
	if err != nil {
		return 0
	}

	return int32(val)
}

func parseInt64(elem string) int64 {
	val, err := strconv.ParseInt(sanitize(elem), 10, 64)
	if err != nil {
		return 0
	}

	return int64(val)
}

func parseFloat32(elem string) float32 {
	val, err := strconv.ParseFloat(elem, 32)
	if err != nil {
		return 0
	}

	return float32(val)
}

func sanitize(elem string) string {
	arr := strings.Split(elem, ",")

	return strings.Join(arr, "")
}
