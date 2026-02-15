package main

import (
"io"
"log"
"os"
"time"

"realtime-poll/config"
"realtime-poll/internal/api"
"realtime-poll/internal/db"
"realtime-poll/internal/middleware"

"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"


)

func main() {


// Load env
err := config.LoadEnv()
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

// Create router
r := gin.New()

// Core middleware
r.Use(gin.Recovery())
r.Use(middleware.RequestLogger())

// =========================
// CORS (GLOBAL MIDDLEWARE)
// =========================
r.Use(cors.New(cors.Config{
	AllowOrigins: []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"https://realtime-poll.clqit.in",
		"https://*.netlify.app",
	},
	AllowMethods: []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
	},
	AllowHeaders: []string{
		"Origin",
		"Content-Type",
		"Authorization",
	},
	ExposeHeaders: []string{
		"Content-Length",
	},
	AllowCredentials: true,
	MaxAge: 12 * time.Hour,
}))

// Register routes AFTER middleware
api.RegisterRoutes(r)

port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}

log.Println("Server running on : http://localhost:" + port)
r.Run(":" + port)

}
