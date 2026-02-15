package middleware

import (
	// "net/http"
	// "time"


	"github.com/gin-gonic/gin"

	"realtime-poll/internal/repository"
	"realtime-poll/internal/constants"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("session_id")
		if err != nil {
			println("NO COOKIE")
			c.AbortWithStatusJSON(401, gin.H{"error": "not logged in"})
			return
		}

		println("COOKIE RECEIVED:", cookie)

		session, err := repository.GetSession(cookie)
		if err != nil || session == nil {
			println("INVALID SESSION")
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid session"})
			return
		}

		println("SESSION USER:", session.UserID)

		c.Set(constants.CtxUserID, session.UserID)
		c.Next()
	}
}
