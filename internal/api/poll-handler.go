
package api

import (
	"net/http"
	"time"
	"strconv"
	"strings"

	"realtime-poll/internal/dto"
	"realtime-poll/internal/services"
	"realtime-poll/internal/constants"
	"realtime-poll/internal/utils"

	"github.com/gin-gonic/gin"
)
func PatchPoll(c *gin.Context) {

	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	// 1️⃣ parse json
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"issue": "invalid json"})
		return
	}

	// 2️⃣ flatten nested json -> mongo dot notation
	flat := map[string]interface{}{}
	utils.Flatten("", body, flat)

	// 3️⃣ protect immutable / system fields
	protected := []string{
		"poll_id",
		"owner_id",

		"meta.created_at",
		"meta.updated_at",
		"meta.total_votes",
		"meta.total_views",
		"meta.last_vote",
		"meta.expires_at",

		"state.version",
		"state.is_deleted",

		"analytics.track_views",
		"analytics.track_voters",
	}

	for _, p := range protected {
		delete(flat, p)
	}

	// prevent vote tampering (no manual vote injection)
	for key := range flat {
		if strings.Contains(key, "votes") {
			delete(flat, key)
		}
	}

	// 4️⃣ always update timestamp
	flat["meta.updated_at"] = time.Now()

	// 5️⃣ call service with FLAT map
	service := services.PollEditService{}
	err := service.PatchPoll(c.Request.Context(), userID, pollID, flat)

	if err != nil {
		c.JSON(400, gin.H{"issue": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "updated"})
}

func PutPoll(c *gin.Context) {

	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"issue": "invalid json"})
		return
	}

	update := map[string]interface{}{}

	// allow only full block replacements
	allowed := []string{
		"content",
		"access",
		"vote",
		"behavior",
		"distribution",
	}

	for _, field := range allowed {
		if val, ok := body[field]; ok {
			update[field] = val
		}
	}

	if len(update) == 0 {
		c.JSON(400, gin.H{"issue": "nothing to update"})
		return
	}

	update["meta.updated_at"] = time.Now()

	service := services.PollEditService{}
	err := service.PutPoll(c.Request.Context(), userID, pollID, update)

	if err != nil {
		c.JSON(400, gin.H{"issue": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "replaced"})
}


func DeletePoll(c *gin.Context) {

	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	force := c.Query("force") == "1"

	service := services.PollEditService{}
	err := service.DeletePoll(c.Request.Context(), userID, pollID, force)

	if err != nil {
		c.JSON(400, gin.H{"issue": err.Error()})
		return
	}

	if force {
		c.JSON(200, gin.H{"status": "permanently deleted"})
		return
	}

	c.JSON(200, gin.H{"status": "soft deleted"})
}


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

func GetPoll(c *gin.Context) {

	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	service := services.PollSingleService{}

	poll, err := service.GetPoll(
		c.Request.Context(),
		userID,
		pollID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"issue": err.Error()})
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
