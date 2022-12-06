package main

import (
	"__elastic/pkg/common/db"
	"__elastic/pkg/common/models"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Shape struct {
	Poster       string
	Title        string
	ReleaseYear  int32
	Certificate  string
	Runtime      string
	Genre        string
	IMDBRating   float32
	Overview     string
	MetaScore    int32
	Director     string
	Votes        int64
	GrossRevenue int64
}

func main() {
	db := db.Init()

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

	db.Transaction(func(tx *gorm.DB) error {
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
