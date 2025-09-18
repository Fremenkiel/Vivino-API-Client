package models

type Wine struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	SeoName         string `json:"seo_name"`
	TypeId          int    `json:"type_id"`
	VintageType     int    `json:"vintage_type"`
	IsNatural       bool   `json:"is_natural"`
	Region          any    `json:"region"`
	VintageMask     any    `json:"vintage_mask"`
	ReviewStatus    int    `json:"review_status"`
	Winery          any    `json:"winery"`
	Style           any    `json:"style"`
	HasValidRatings bool   `json:"has_valid_ratings"`
}
