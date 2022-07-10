package models

import "time"

type Ad struct {
	ID      int    `json:"ad_id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	//Cat       Category  `json:"cat"`
	Category  string    `json:"category" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
