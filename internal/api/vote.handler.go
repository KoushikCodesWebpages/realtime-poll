package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"realtime-poll/internal/apperror"
	"realtime-poll/internal/constants"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/services"
)

type VoteReq struct {
	PollID   string `json:"poll_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

func CastVote(c *gin.Context) {

	var req VoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			apperror.New(apperror.TokenMalformed, "invalid request body"))
		return
	}

	ctx := c.Request.Context()

	// -------------------- identity --------------------
	userID := c.GetString(constants.CtxUserID)

	sessionID, _ := c.Cookie("session_id") // guest allowed
	ip := c.ClientIP()

	// -------------------- load poll --------------------
	poll, err := repository.GetPollByID(ctx, req.PollID)
	if err != nil || poll == nil {
		c.JSON(http.StatusNotFound,
			apperror.New(apperror.PollNotFound, "poll not found"))
		return
	}

	voteService := services.VoteService{}

	// -------------------- cast vote --------------------
	err = voteService.CastVote(
		ctx,
		req.PollID,
		req.OptionID,
		userID,
		sessionID,
		ip,
	)
	if err != nil {
		c.JSON(http.StatusForbidden,
			apperror.New(apperror.PollNotAllowed, err.Error()))
		return
	}

	// -------------------- response --------------------
	// For realtime polls → WS will broadcast results
	if poll.Behavior.ShowLiveResults {
		c.JSON(http.StatusOK, gin.H{
			"status": "vote_received",
		})
		return
	}

	// Non realtime → return results if allowed
	results, err := voteService.GetResults(ctx, req.PollID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "vote_recorded"})
		return
	}

	canViewResults := !poll.Vote.HideResults

	c.JSON(http.StatusOK, gin.H{
		"status":           "vote_recorded",
		"results_visible":  canViewResults,
		"results":          results,
	})
}
