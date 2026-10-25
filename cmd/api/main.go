package main

import (
	"books-api/app/configs"
	"books-api/app/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("load .env: ", err)
	}
	if err := configs.OpenDatabase(); err != nil {
		log.Fatal(err)
	}
	defer configs.DB.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be set and contain at least 32 characters")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()
	api := router.Group("/api")

	routes.RegisterAuthRoutes(api)
	routes.RegisterCategoryRoutes(api)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("start API server: ", err)
	}
}
