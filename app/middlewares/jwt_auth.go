package middlewares

import (
	"net/http"
	"os"
	"strings"

	"books-api/app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		parts := strings.Fields(context.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helpers.Respond(context, http.StatusUnauthorized, "valid bearer token is required", nil)
			context.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if len(secret) < 32 {
			helpers.Respond(context, http.StatusInternalServerError, "authentication is not configured", nil)
			context.Abort()
			return
		}
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			helpers.Respond(context, http.StatusUnauthorized, "invalid or expired bearer token", nil)
			context.Abort()
			return
		}
		context.Next()
	}
}
