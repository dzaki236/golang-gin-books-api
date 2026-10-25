package routes

import (
	"books-api/app/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router gin.IRouter) {
	auth := router.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
}
