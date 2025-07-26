package models

import "time"

// Banner represents a banner entity
type Banner struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// ClickStat represents a click statistics entry for a banner
type ClickStat struct {
	BannerID  int       `json:"banner_id" db:"banner_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Count     int       `json:"count" db:"count"`
}

// StatsRequest represents the request for statistics API
type StatsRequest struct {
	From string `json:"from" example:"2024-01-01T00:00:00"`
	To   string `json:"to" example:"2024-01-02T00:00:00"`
}

// StatsResponse represents the response for statistics API
type StatsResponse struct {
	Stats []StatEntry `json:"stats"`
}

// StatEntry represents a statistics entry in the API response
type StatEntry struct {
	Timestamp string `json:"ts" example:"2024-01-01T00:00:00"`
	Value     int    `json:"v" example:"5"`
}
