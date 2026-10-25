package structs

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
