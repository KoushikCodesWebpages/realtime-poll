package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"realtime-poll/internal/repository"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
			return
		}

		session, err := repository.GetSession(cookie)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
