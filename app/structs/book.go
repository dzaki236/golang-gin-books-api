package structs

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ReleaseYear int       `json:"release_year"`
	Price       int       `json:"price"`
	TotalPage   int       `json:"total_page"`
	Thickness   string    `json:"thickness"`
	CategoryID  int       `json:"category_id"`
	Category    Category  `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required,max=2048"`
	ReleaseYear int    `json:"release_year" binding:"required,min=1980,max=2024"`
	Price       int    `json:"price" binding:"required,gte=0"`
	TotalPage   int    `json:"total_page" binding:"required,gt=0"`
	CategoryID  int    `json:"category_id" binding:"required,gt=0"`
}
