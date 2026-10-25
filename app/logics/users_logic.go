package logics

import (
	"context"
	"database/sql"
	"errors"

	"books-api/app/configs"
	"books-api/app/helpers"
	"books-api/app/structs"
)

func CreateUser(ctx context.Context, username, passwordHash string) (structs.User, error) {
	var user structs.User
	err := configs.DB.QueryRowContext(ctx,
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at`,
		username, passwordHash,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return structs.User{}, helpers.ErrUsernameTaken
	}
	return user, nil
}

func FindUserByUsername(ctx context.Context, username string) (structs.User, error) {
	var user structs.User
	err := configs.DB.QueryRowContext(ctx,
		`SELECT id, username, password, created_at FROM users WHERE username = $1`,
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
