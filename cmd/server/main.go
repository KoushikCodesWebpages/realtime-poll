package main

import (
	"io"
	"log"
	"os"
	"strings"
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

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// =========================
	// CORS (FIXED)
	// =========================
	allowedOrigins := config.LoadCorsOrigins()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {

			// exact origins from .env
			if allowedOrigins[origin] {
				return true
			}

			// allow all netlify preview deployments
			if strings.HasSuffix(origin, ".netlify.app") {
				return true
			}

			return false
		},

		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Routes
	api.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on : http://localhost:" + port)
	r.Run(":" + port)
}
