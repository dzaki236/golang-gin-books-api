package routes

import (
	"books-api/app/handlers"
	"books-api/app/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/logout", middlewares.JWTAuth(), handlers.Logout)
}
