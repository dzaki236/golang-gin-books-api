package main

import (
	"books-api/app/configs"
	"books-api/app/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("gagal memuat berkas .env: ", err)
	}
	if err := configs.OpenDatabase(); err != nil {
		log.Fatal(err)
	}
	defer configs.DB.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET wajib diatur dan berisi setidaknya 32 karakter")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "api siap digunakan",
		})
	})
	routes.RegisterAuthRoutes(api)
	routes.RegisterCategoryRoutes(api)
	routes.RegisterBookRoutes(api)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("gagal menjalankan server API: ", err)
	}
}
