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

	poll, _ := repository.GetPollByIDRaw(c.Request.Context(), pollID)
	if poll == nil || poll.OwnerID != userID {
		c.JSON(403, gin.H{"issue": "not allowed"})
		return
	}

	var body struct {
		Mode    string `json:"mode"` // infinite | timed
		Minutes int64  `json:"minutes"`
	}

	c.ShouldBindJSON(&body)

	if body.Mode == "" {
		body.Mode = "infinite"
	}

	token, err := services.GenerateShareToken(pollID, body.Mode, body.Minutes)
	if err != nil {
		c.JSON(500, gin.H{"issue": "failed to create link"})
		return
	}

	c.JSON(200, gin.H{
		"link": "/b1/poll/share?token=" + token,
	})
}

func ViewSharedPoll(c *gin.Context) {

	token := c.Query("token")

	claims, err := services.VerifyShareToken(token)
	if err != nil {
		c.JSON(401, gin.H{"issue": err.Error()})
		return
	}

	poll, _ := repository.GetPollByID(c.Request.Context(), claims.PollID)
	if poll == nil {
		c.JSON(404, gin.H{"issue": "poll not found"})
		return
	}

	c.JSON(200, poll)
}