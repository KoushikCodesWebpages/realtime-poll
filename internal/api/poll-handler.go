package api

import (
	// "net/http"
	"realtime-poll/internal/services"

	"github.com/gin-gonic/gin"
)

type CreatePollReq struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

// func CreatePoll(c *gin.Context) {

// 	userID := c.GetString("user_id") // set by auth middleware

// 	var req models.CreatePollReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{"error": "invalid body"})
// 		return
// 	}

// 	poll, err := services.CreatePoll(userID, req.Question, req.Options)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, poll)
// }

func CreatePoll(c *gin.Context) {

	userID := c.GetString("user_id") // set by middleware

	c.JSON(200, gin.H{
		"message": "poll creation allowed",
		"owner":   userID,
	})
}

type VoteReq struct {
	OptionID string `json:"option_id"`
}

func Vote(c *gin.Context) {
	pollID := c.Param("id")

	var req VoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad"})
		return
	}

	err := services.Vote(pollID, req.OptionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
}
