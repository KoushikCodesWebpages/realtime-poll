package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"realtime-poll/internal/apperror"
	"realtime-poll/internal/constants"
	"realtime-poll/internal/middleware"
	"realtime-poll/internal/services"
)

/* ---------------- Request ---------------- */

type VoteReq struct {
	PollID   string `json:"poll_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

/* ---------------- Handler ---------------- */

func CastVoteHandler(c *gin.Context) {

	var req VoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.BadRequest("Invalid request body"))
		return
	}

	ctx := c.Request.Context()

	userID := c.GetString(constants.CtxUserID)
	sessionID, _ := c.Cookie("session_id")
	ip := c.ClientIP()

	voteService := services.VoteService{}

	result, err := voteService.CastVote(
		ctx,
		req.PollID,
		req.OptionID,
		userID,
		sessionID,
		ip,
	)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"viewer": gin.H{
			"can_vote":        false,
			"selected_option": result.SelectedOption,
			"already_voted":   result.AlreadyVoted,
			"changed":         result.Changed,
		},
		"version": result.Version,
	})
}
