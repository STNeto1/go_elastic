package models

type Movie struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	Poster       string  `json:"poster"`
	Title        string  `json:"title"`
	ReleaseYear  int32   `json:"release_year"`
	Certificate  string  `json:"certificate"`
	Runtime      string  `json:"runime"`
	Genre        string  `json:"genre"`
	IMDBRating   float32 `json:"imdb_rating"`
	Overview     string  `json:"overview"`
	MetaScore    int32   `json:"meta_score"`
	Director     string  `json:"director"`
	Votes        int64   `json:"votes"`
	GrossRevenue int64   `json:"gross_revenue"`
}
