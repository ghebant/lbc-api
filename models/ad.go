package models

import "time"

type Ad struct {
	ID      int    `json:"ad_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	//Category  Category  `json:"category"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
