package models

type Statistics struct {
	Status         string  `json:"status"`
	RatingsCount   int     `json:"ratings_count"`
	RatingsAverage float32 `json:"ratings_average"`
	LabelsCount    int     `json:"labels_count"`
	ReviewsCount   int     `json:"reviews_count"`
}
