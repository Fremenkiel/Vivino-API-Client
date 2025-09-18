package models

type Vintage struct {
	Id                     int        `json:"id"`
	SeoName                string     `json:"seo_name"`
	Name                   string     `json:"name"`
	Statistics             Statistics `json:"statistics"`
	OrganicCertificationId any        `json:"organic_certification_id"`
	CertifiedBiodynamic    any        `json:"certified_biodynamic"`
	Image                  any        `json:"image"`
	Wine                   Wine       `json:"wine"`
	Year                   string     `json:"year"`
	Grapes                 any        `json:"grapes"`
	HasHalidRatings        bool       `json:"has_valid_ratings"`
}
