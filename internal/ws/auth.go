package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"realtime-poll/internal/utils"
	"realtime-poll/internal/repository"
)

func resolveSession(c *gin.Context, pollID string) *Session {

	token := c.Query("token")

	sub, _ := utils.ParseWSToken(token)

	var userID *string
	if sub != "" {
		userID = &sub
	}

	sess := &Session{
		ConnID:   uuid.NewString(),
		TokenSub: sub,
		UserID:   userID,
		PollID:   pollID,
	}

	// determine ownership ONCE
	if userID != nil {
		poll, _ := repository.GetPollByID(c.Request.Context(), pollID)
		if poll != nil && poll.OwnerID == *userID {
			sess.IsOwner = true
		}
	}

	return sess
}