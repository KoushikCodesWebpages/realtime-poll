package middleware

import (
	"github.com/gin-gonic/gin"
	"realtime-poll/internal/models"
)

func GetUser(c *gin.Context) (*models.User, bool) {
	u, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := u.(*models.User)
	return user, ok
}
