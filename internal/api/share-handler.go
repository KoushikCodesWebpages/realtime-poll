// package api

// import (
// 	"net/http"
// 	"os"
// 	"time"

// 	"realtime-poll/internal/services"
// 	"realtime-poll/internal/constants"
// 	"realtime-poll/internal/repository"
// 	"realtime-poll/internal/apperror"
// 	"realtime-poll/internal/utils"

// 	"github.com/gin-gonic/gin"
// )

package api

import (
	"net/http"

	"time"
	"github.com/gin-gonic/gin"


	"realtime-poll/internal/constants"
	"realtime-poll/internal/middleware"

	"realtime-poll/internal/services"

)
var shareService = services.NewShareService()

func GenerateShareLink(c *gin.Context) {

	ctx := c.Request.Context()
	pollID := c.Param("poll_id")
	userID := c.GetString(constants.CtxUserID)

	token, mode, err := shareService.GenerateShareLink(ctx, userID, pollID)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"share_url": "/share?token=" + token,
		"token":     token,
		"mode":      mode,
	})
}

func ViewSharedPoll(c *gin.Context) {

	ctx := c.Request.Context()
	token := c.Query("token")

	sessionID, _ := c.Cookie("session_id")

	result, err := shareService.ViewSharedPoll(
		ctx,
		token,
		c.ClientIP(),
		sessionID,
	)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	// if new session created → set cookie
	if sessionID == "" && result.SessionID != "" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session_id",
			Value:    result.SessionID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		})
	}

	c.JSON(http.StatusOK, result)
}
