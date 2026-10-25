package logics

import (
	"database/sql"
	"errors"
	"strings"

	"books-api/app/configs"
	"books-api/app/helpers"
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
)

func ListCategories(ctx *gin.Context) ([]structs.Category, error) {
	rows, err := configs.DB.QueryContext(ctx.Request.Context(), `SELECT id, name, created_at FROM categories ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]structs.Category, 0)
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func FindCategory(ctx *gin.Context, id int) (structs.Category, error) {
	var category structs.Category
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`SELECT id, name, created_at FROM categories WHERE id = $1`, id,
	).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Category{}, helpers.ErrCategoryNotFound
	}
	return category, err
}

func CreateCategory(ctx *gin.Context, name, username string) (structs.Category, error) {
	var category structs.Category
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`INSERT INTO categories (name, created_by) VALUES ($1, $2) RETURNING id, name, created_at`, name, username,
	).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "categories_name_key") {
			return structs.Category{}, helpers.ErrCategoryExists
		}
		return structs.Category{}, err
	}
	return category, nil
}

func UpdateCategory(ctx *gin.Context, id int, name, username string) (structs.Category, error) {
	var category structs.Category
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`UPDATE categories SET name = $2, modified_at = CURRENT_TIMESTAMP, modified_by = $3 WHERE id = $1 RETURNING id, name, created_at`,
		id, name, username,
	).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Category{}, helpers.ErrCategoryNotFound
	}
	if err != nil {
		if strings.Contains(err.Error(), "categories_name_key") {
			return structs.Category{}, helpers.ErrCategoryExists
		}
		return structs.Category{}, err
	}
	return category, nil
}

func DeleteCategory(ctx *gin.Context, id int) error {
	result, err := configs.DB.ExecContext(ctx.Request.Context(), `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return helpers.ErrCategoryNotFound
	}
	return nil
}
