
package api

import (
	"net/http"

	"realtime-poll/internal/dto"
	"realtime-poll/internal/services"
	"realtime-poll/internal/constants"

	"github.com/gin-gonic/gin"
)

func CreatePoll(c *gin.Context) {

	var req dto.CreatePollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"issue": err.Error()})
		return
	}

	userID := c.GetString(constants.CtxUserID)
	email := c.GetString("auth_email")

	service := services.PollCreateService{}

	poll, err := service.CreatePoll(c.Request.Context(), userID, email, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"issue": err.Error()})
		return
	}

	c.JSON(http.StatusOK, poll)
}


// func CreatePoll(c *gin.Context) {

// 	userID := c.GetString("user_id") // set by middleware

// 	c.JSON(200, gin.H{
// 		"message": "poll creation allowed",
// 		"owner":   userID,
// 	})
// }
