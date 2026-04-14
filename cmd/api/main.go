package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/OtavMacedo/url-shortener-golang/internal/cache"
	"github.com/OtavMacedo/url-shortener-golang/internal/controller"
	"github.com/OtavMacedo/url-shortener-golang/internal/infra"
	"github.com/OtavMacedo/url-shortener-golang/internal/middleware"
	"github.com/OtavMacedo/url-shortener-golang/internal/ratelimiter"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	jwtExpirationHours := os.Getenv("JWT_EXPIRATION_HOURS")
	if jwtExpirationHours == "" {
		log.Fatal("JWT_EXPIRATION_HOURS environment variable is required")
	}
	jwtExpirationHoursInt, err := strconv.Atoi(jwtExpirationHours)
	if err != nil {
		log.Fatal("Invalid JWT_SECRET environment variable")
	}
	jwtExpiration := time.Duration(jwtExpirationHoursInt) * time.Hour

	pool, err := infra.ConnectDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	redisClient, err := infra.ConnectRedis(redisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	authService := service.NewAuthService(jwtSecret, jwtExpiration)
	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, authService)
	userController := controller.NewUserController(userService)
	urlRepository := repository.NewUrlRepository(pool)
	redisCache := cache.NewRedisCache(redisClient)
	urlService := service.NewUrlService(urlRepository, redisCache)
	urlController := controller.NewUrlController(urlService)
	rateLimiter := ratelimiter.NewRateLimiter(redisClient, 5, 1*time.Minute)

	router := gin.Default()
	router.POST("/users", userController.Create)
	router.POST("/login", userController.Login)
	router.GET("/:slug", urlController.Redirect)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))

	protected.POST("/urls", middleware.RateLimiterMiddleware(rateLimiter), urlController.Create)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start api: %v", err)
	}
}
