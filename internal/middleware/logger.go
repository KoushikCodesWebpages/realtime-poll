package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	debug := os.Getenv("DEBUG") == "true"

	return func(c *gin.Context) {

		// skip logging if debug disabled
		if !debug {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		latency := time.Since(start)

		log.Printf(
			"%s | %d | %v | %s %s",
			c.ClientIP(),
			c.Writer.Status(),
			latency,
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}
