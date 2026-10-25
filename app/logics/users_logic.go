package logics

import (
	"database/sql"
	"errors"

	"books-api/app/configs"
	"books-api/app/helpers"
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context, username, passwordHash string) (structs.User, error) {
	var user structs.User
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at`,
		username, passwordHash,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return structs.User{}, helpers.ErrUsernameTaken
	}
	return user, nil
}

func FindUserByUsername(ctx *gin.Context, username string) (structs.User, error) {
	var user structs.User
	err := configs.DB.QueryRowContext(ctx.Request.Context(),
		`SELECT id, username, password, created_at FROM users WHERE username = $1 LIMIT 1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.User{}, helpers.ErrUserNotFound
	}
	if err != nil {
		return structs.User{}, err
	}
	return user, nil
}
