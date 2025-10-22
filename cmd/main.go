package main

import (
	"database/sql"

	"os"
	"rest-service/internal/handlers"
	"rest-service/internal/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	docs "rest-service/docs"

	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"client": c.ClientIP(),
		}).Info("Incoming request")

		c.Next()

		status := c.Writer.Status()
		log.WithFields(log.Fields{
			"status": status,
		}).Info("Request completed")
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
		log.Infof("DATABASE_URL: %s", os.Getenv("DATABASE_URL"))
	}
	log.Infof("DATABASE_URL: %s", os.Getenv("DATABASE_URL"))

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	repo := repository.NewPostgresSubscriptionRepository(db)
	handler := handlers.NewSubscriptionHandler(repo)

	r := gin.Default()
	r.Use(LoggerMiddleware())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = "/"

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to rest-service API"})
	})

	r.POST("/subscriptions", handler.Create)
	r.GET("/subscriptions", handler.GetAll)
	r.GET("/subscriptions/:id", handler.GetByID)
	r.PUT("/subscriptions/:id", handler.Update)
	r.DELETE("/subscriptions/:id", handler.Delete)
	r.GET("/subscriptions/sum", handler.GetSum)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
