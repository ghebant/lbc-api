package models

import (
	"ghebant/lbc-api/internal/constants"
	"time"
)

type Ad struct {
	ID         int                `json:"ad_id"`
	Title      string             `json:"title" binding:"required"`
	Content    string             `json:"content" binding:"required"`
	Category   constants.Category `json:"category" binding:"required"`
	Automobile *Automobile        `json:"automobile,omitempty"`
	RealEstate *RealEstate        `json:"real_estate,omitempty"`
	Job        *Job               `json:"job,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
