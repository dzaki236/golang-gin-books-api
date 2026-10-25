package routes

import (
	"books-api/app/handlers"
	"books-api/app/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(api *gin.RouterGroup) {
	books := api.Group("/books", middlewares.JWTAuth())
	books.GET("", handlers.ListBooks)
	books.POST("", handlers.CreateBook)
	books.GET("/:id", handlers.GetBook)
	books.PUT("/:id", handlers.UpdateBook)
	books.DELETE("/:id", handlers.DeleteBook)
}
