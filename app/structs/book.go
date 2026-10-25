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
	CreatedAt   time.Time `json:"created_at"`
}

type BookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year"`
	Price       int    `json:"price"`
	TotalPage   int    `json:"total_page"`
	CategoryID  int    `json:"category_id"`
}
