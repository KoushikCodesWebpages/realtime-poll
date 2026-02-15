package api

import (

	"realtime-poll/internal/services"
	"realtime-poll/internal/constants"
	"realtime-poll/internal/repository"


	"github.com/gin-gonic/gin"
)

func GenerateShareLink(c *gin.Context) {

	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	// TODO: verify ownership later
	if userID == "" {
		c.JSON(401, gin.H{"issue": "unauthorized"})
		return
	}

	token, err := services.GenerateShareToken(pollID)
	if err != nil {
		c.JSON(500, gin.H{"issue": "failed to generate"})
		return
	}

	link := "/b1/poll/share?token=" + token

	c.JSON(200, gin.H{"link": link})
}

func ViewSharedPoll(c *gin.Context) {

	token := c.Query("token")

	claims, err := services.VerifyShareToken(token)
	if err != nil {
		c.JSON(401, gin.H{"issue": err.Error()})
		return
	}

	poll, err := repository.GetPollByID(c.Request.Context(), claims.PollID)
	if err != nil {
		c.JSON(404, gin.H{"issue": "poll not found"})
		return
	}

	c.JSON(200, poll)
}
