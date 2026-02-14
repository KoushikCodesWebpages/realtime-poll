package main

import (
	"log"
	"os"
	"io"

	"realtime-poll/internal/api"
	"realtime-poll/internal/db"

	"realtime-poll/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	// DB connect
	db.Connect()

	debug := os.Getenv("DEBUG") == "true"

	if !debug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// Register all routes
	api.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on : http://localhost:" + port)
	r.Run(":" + port)
}
