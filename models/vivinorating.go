package models

import "time"

type VivinoRating struct {
	Id                int       `json:"id"`
	Rating            float32   `json:"rating"`
	Note              string    `json:"note"`
	Language          any       `json:"language"`
	CreatedAt         time.Time `json:"created_at"`
	Aggregated        bool      `json:"aggregated"`
	User              any       `json:"user"`
	Vintage           Vintage   `json:"vintage"`
	Activity          any       `json:"activity"`
	FlavorWordMatches any       `json:"flavor_word_matches"`
	TaggedNote        string    `json:"tagged_note"`
}
