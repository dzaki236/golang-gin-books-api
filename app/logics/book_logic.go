package logics

import (
	"context"
	"database/sql"
	"errors"

	"books-api/app/configs"
	"books-api/app/helpers"
	"books-api/app/structs"
)

func ListBooks(ctx context.Context) ([]structs.Book, error) {
	rows, err := configs.DB.QueryContext(ctx, `SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]structs.Book, 0)
	for rows.Next() {
		var book structs.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func FindBook(ctx context.Context, id int) (structs.Book, error) {
	var book structs.Book
	err := configs.DB.QueryRowContext(ctx,
		`SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at FROM books WHERE id = $1`, id,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Book{}, helpers.ErrBookNotFound
	}
	return book, err
}

func CreateBook(ctx context.Context, request structs.BookRequest) (structs.Book, error) {
	var book structs.Book
	var thickness string
	if request.TotalPage > 100 {
		thickness = "tebal"
	}
	thickness = "tipis"
	err := configs.DB.QueryRowContext(ctx,
		`INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at`,
		request.Title, request.Description, request.ImageURL, request.ReleaseYear, request.Price, request.TotalPage, thickness, request.CategoryID,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt)
	return book, err
}

func UpdateBook(ctx context.Context, id int, request structs.BookRequest) (structs.Book, error) {
	var book structs.Book

	var thickness string
	if request.TotalPage > 100 {
		thickness = "tebal"
	}
	thickness = "tipis"
	err := configs.DB.QueryRowContext(ctx,
		`UPDATE books SET title = $2, description = $3, image_url = $4, release_year = $5, price = $6, total_page = $7, thickness = $8, category_id = $9, modified_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at`,
		id, request.Title, request.Description, request.ImageURL, request.ReleaseYear, request.Price, request.TotalPage, thickness, request.CategoryID,
	).Scan(&book.ID, &book.Title, &book.Description, &book.ImageURL, &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness, &book.CategoryID, &book.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Book{}, helpers.ErrBookNotFound
	}
	return book, err
}

func DeleteBook(ctx context.Context, id int) error {
	result, err := configs.DB.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
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
