package models

import "time"

type Automobile struct {
	ID        int       `json:"id" db:"automobile_id"`
	AdId      int       `json:"ad_id"`
	Brand     string    `json:"brand" binding:"required""`
	Model     string    `json:"model" binding:"required""`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
