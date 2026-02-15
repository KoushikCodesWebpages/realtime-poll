
package api

import (
	"net/http"
	"time"
	"strconv"

	"realtime-poll/internal/dto"
	"realtime-poll/internal/services"
	"realtime-poll/internal/constants"

	"github.com/gin-gonic/gin"
)
func parsePollFilter(c *gin.Context) dto.PollFilter {

	limit := int64(10)
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = int64(v)
		}
	}

	var hasVotes *bool
	if hv := c.Query("has_votes"); hv != "" {
		val := hv == "true"
		hasVotes = &val
	}

	var dateFrom *time.Time
	if df := c.Query("date_from"); df != "" {
		if t, err := time.Parse(time.RFC3339, df); err == nil {
			dateFrom = &t
		}
	}

	var dateTo *time.Time
	if dt := c.Query("date_to"); dt != "" {
		if t, err := time.Parse(time.RFC3339, dt); err == nil {
			dateTo = &t
		}
	}

	return dto.PollFilter{
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		Visibility: c.Query("visibility"),
		HasVotes:   hasVotes,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Sort:       c.DefaultQuery("sort", "newest"),
		Limit:      limit,
	}
}
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

func GetMyPolls(c *gin.Context) {

	userID := c.GetString(constants.CtxUserID)

	// 1️⃣ parse filters
	filter := parsePollFilter(c)

	// 2️⃣ cursor + direction
	cursorStr := c.Query("cursor")
	direction := c.DefaultQuery("dir", "next")

	var cursor *time.Time
	if cursorStr != "" {
		t, err := time.Parse(time.RFC3339, cursorStr)
		if err == nil {
			cursor = &t
		}
	}

	// 3️⃣ preserve filters for pagination links
	query := c.Request.URL.Query()
	query.Del("cursor")
	query.Del("dir")

	service := services.PollGetService{}

	resp, err := service.GetMyPollsPaginated(
		c.Request.Context(),
		userID,
		filter,
		query,
		cursor,
		direction,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"issue": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// func CreatePoll(c *gin.Context) {

// 	userID := c.GetString("user_id") // set by middleware

// 	c.JSON(200, gin.H{
// 		"message": "poll creation allowed",
// 		"owner":   userID,
// 	})
// }
