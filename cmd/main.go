package main

import (
	"database/sql"
	"log"
	"os"
	"rest-service/internal/handlers"
	"rest-service/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewPostgresSubscriptionRepository(db)
	handler := handlers.NewSubscriptionHandler(repo)

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/subscriptions", handler.Create)
	r.GET("/subscriptions", handler.GetAll)
	r.GET("/subscriptions/:id", handler.GetByID)
	r.PUT("/subscriptions/:id", handler.Update)
	r.DELETE("/subscriptions/:id", handler.Delete)
	r.GET("/subscriptions/sum", handler.GetSum)

	r.Run(":8080")
}
