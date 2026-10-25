package logics

import (
	"database/sql"
	"errors"

	"books-api/app/configs"
	"books-api/app/helpers"
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
)

func ListBooks(ctx *gin.Context) ([]structs.Book, error) {
	rows, err := configs.DB.QueryContext(ctx.Request.Context(), `SELECT book.id, book.title, book.description, book.image_url, book.release_year, book.price, book.total_page, book.thickness, book.category_id, book.created_at, category.id, category.name, category.created_at FROM books book JOIN categories category ON category.id = book.category_id ORDER BY book.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]structs.Book, 0)
	for rows.Next() {
		var book structs.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.Category.ID, &book.Category.Name, &book.Category.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func ListBooksByCategory(ctx *gin.Context, categoryID int) ([]structs.Book, error) {
	rows, err := configs.DB.QueryContext(ctx.Request.Context(), `SELECT book.id, book.title, book.description, book.image_url, book.release_year, book.price, book.total_page, book.thickness, book.category_id, book.created_at, category.id, category.name, category.created_at FROM books book JOIN categories category ON category.id = book.category_id WHERE category.id = $1 ORDER BY book.created_at DESC`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]structs.Book, 0)
	for rows.Next() {
		var book structs.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.Category.ID, &book.Category.Name, &book.Category.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func FindBook(ctx *gin.Context, id int) (structs.Book, error) {
	var book structs.Book
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`SELECT book.id, book.title, book.description, book.image_url, book.release_year, book.price, book.total_page, book.thickness, book.category_id, book.created_at, category.id, category.name, category.created_at FROM books book JOIN categories category ON category.id = book.category_id WHERE book.id = $1`, id,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt, &book.Category.ID, &book.Category.Name, &book.Category.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Book{}, helpers.ErrBookNotFound
	}
	return book, err
}

func CreateBook(ctx *gin.Context, request structs.BookRequest, username string) (structs.Book, error) {
	var book structs.Book
	var thickness string
	if request.TotalPage > 100 {
		thickness = "tebal"
	}
	thickness = "tipis"
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at`,
		request.Title, request.Description, request.ImageURL, request.ReleaseYear, request.Price, request.TotalPage, thickness, request.CategoryID, username,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt)
	return book, err
}

func UpdateBook(ctx *gin.Context, id int, request structs.BookRequest, username string) (structs.Book, error) {
	var book structs.Book

	var thickness string
	if request.TotalPage > 100 {
		thickness = "tebal"
	}
	thickness = "tipis"
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`UPDATE books SET title = $2, description = $3, image_url = $4, release_year = $5, price = $6, total_page = $7, thickness = $8, category_id = $9, modified_at = CURRENT_TIMESTAMP, modified_by = $10 WHERE id = $1 RETURNING id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at`,
		id, request.Title, request.Description, request.ImageURL, request.ReleaseYear, request.Price, request.TotalPage, thickness, request.CategoryID, username,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Book{}, helpers.ErrBookNotFound
	}
	return book, err
}

func DeleteBook(ctx *gin.Context, id int) error {
	result, err := configs.DB.ExecContext(ctx.Request.Context(), `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return helpers.ErrBookNotFound
	}
	return nil
}
