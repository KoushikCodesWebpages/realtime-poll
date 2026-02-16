package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"realtime-poll/internal/middleware"
	"realtime-poll/internal/services"
)

var snapshotService = services.NewSnapshotService()

func GetPollSnapshot(c *gin.Context) {

	pollID := c.Param("poll_id")

	data, err := snapshotService.GetSnapshot(c.Request.Context(), pollID)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
