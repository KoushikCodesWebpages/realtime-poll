package api

import (

	"realtime-poll/internal/services"
	"realtime-poll/internal/constants"
	"net/http"


	"github.com/gin-gonic/gin"
)
func CastVote(c *gin.Context) {

	var body struct {
		PollID   string `json:"poll_id"`
		OptionID string `json:"option_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"issue": "invalid body"})
		return
	}

	// user may or may not exist
	userID := c.GetString(constants.CtxUserID)

	sessionID, _ := c.Cookie("session_id")
	ip := c.ClientIP()

	service := services.VoteService{}
	err := service.CastVote(c.Request.Context(), body.PollID, body.OptionID, userID, sessionID, ip)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"issue": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "voted"})
}
