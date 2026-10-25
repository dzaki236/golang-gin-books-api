package structs

import "github.com/golang-jwt/jwt/v5"

type JWTUsersSign struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
