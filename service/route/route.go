package route

import "time"

type Route struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Type        string    `json:"type"`
	Operator    string    `json:"operator"`
	DepartureAt time.Time `json:"departure_at"`
	ArrivalAt   time.Time `json:"arrival_at"`
	Price       float64   `json:"price"`
	Quota       int       `json:"quota"`
	Sold        int       `json:"sold"`
	Available   int       `json:"available" gorm:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchFilter struct {
	Origin      string
	Destination string
	Type        string
	Date        string
}
