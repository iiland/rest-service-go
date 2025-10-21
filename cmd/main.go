package main

import (
	"database/sql"
	"log"
	"os"
	"rest-service/internal/handlers"
	"rest-service/internal/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	// Применяем миграции из папки migrations
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	repo := repository.NewPostgresSubscriptionRepository(db)
	handler := handlers.NewSubscriptionHandler(repo)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/subscriptions", handler.Create)
	router.GET("/subscriptions", handler.GetAll)
	router.GET("/subscriptions/:id", handler.GetByID)
	router.PUT("/subscriptions/:id", handler.Update)
	router.DELETE("/subscriptions/:id", handler.Delete)
	router.GET("/subscriptions/sum", handler.GetSum)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
