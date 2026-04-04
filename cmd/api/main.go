package main

import (
	"log"
	"os"

	"github.com/OtavMacedo/url-shortener-golang/internal/controller"
	"github.com/OtavMacedo/url-shortener-golang/internal/infra"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pool, err := infra.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router := gin.Default()
	router.POST("/users", userController.Create)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start api: %v", err)
	}
}
