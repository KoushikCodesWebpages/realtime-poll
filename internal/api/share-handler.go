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

	incomingSessionID, _ := c.Cookie("session_id")

	result, err := shareService.ViewSharedPoll(
		ctx,
		token,
		c.ClientIP(),
		incomingSessionID,
	)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	isProd := gin.Mode() == gin.ReleaseMode

	// --------------------------------------------------
	// ACCESS BLOCK → CLEAR SESSION COOKIE
	// --------------------------------------------------
	if result.Access == "not_whitelisted" || result.Access == "whitelist_required" {

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // delete cookie
			HttpOnly: true,
		}

		if isProd {
			cookie.Domain = ".clqit.in"
			cookie.Secure = true
			cookie.SameSite = http.SameSiteNoneMode
		} else {
			cookie.Secure = false
			cookie.SameSite = http.SameSiteLaxMode
		}

		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, result)
		return
	}

	// --------------------------------------------------
	// NORMAL FLOW → SYNC SESSION COOKIE
	// --------------------------------------------------
	if result.SessionID != "" && result.SessionID != incomingSessionID {

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    result.SessionID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		}

		if isProd {
			cookie.Domain = ".clqit.in"
			cookie.Secure = true
			cookie.SameSite = http.SameSiteNoneMode
		} else {
			cookie.Secure = false
			cookie.SameSite = http.SameSiteLaxMode
		}

		http.SetCookie(c.Writer, cookie)
	}

	c.JSON(http.StatusOK, result)
}
