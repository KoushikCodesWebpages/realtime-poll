package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func resolveSession(c *gin.Context, pollID string) *Session {

	var userID *string

	if v, ok := c.Get("user_id"); ok {
		id := v.(string)
		userID = &id
	}

	return &Session{
		ConnID:    uuid.NewString(),
		UserID:    userID,
		Anonymous: uuid.NewString(),
		PollID:    pollID,
		Role:      RoleViewer,
		Visibility: ViewLive,
	}
}
