package routes

import (
	"books-api/app/handlers"
	"books-api/app/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	categories := router.Group("/categories", middlewares.JWTAuth())
	categories.GET("", handlers.ListCategories)
	categories.POST("", handlers.CreateCategory)
	categories.GET("/:id/books", handlers.ListCategoryBooks)
	categories.GET("/:id", handlers.GetCategory)
	categories.PUT("/:id", handlers.UpdateCategory)
	categories.DELETE("/:id", handlers.DeleteCategory)
}
