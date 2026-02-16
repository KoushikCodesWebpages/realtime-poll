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
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"realtime-poll/internal/apperror"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/constants"
	"realtime-poll/internal/services"
	"realtime-poll/internal/utils"
)


func GenerateShareLink(c *gin.Context) {

	ctx := c.Request.Context()
	pollID := c.Param("poll_id")

	// must be owner
	userID := c.GetString(constants.CtxUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized,
			apperror.New(apperror.TokenInvalid, "login required"))
		return
	}

	// load poll
	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil {
		c.JSON(http.StatusNotFound,
			apperror.New(apperror.PollNotFound, "poll not found"))
		return
	}

	if poll.OwnerID != userID {
		c.JSON(http.StatusForbidden,
			apperror.New(apperror.PollNotAllowed, "not poll owner"))
		return
	}

	// generate token based on visibility
	duration := 365 * 24 * time.Hour // 1 year share link

	token, err := utils.GenerateShareToken(
		poll.PollID,
		poll.Access.Visibility,
		duration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			apperror.New(apperror.TokenInvalid, "failed to generate link"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_url": "/share?token=" + token,
		"token":     token,
	})
}


func ViewSharedPoll(c *gin.Context) {

	ctx := c.Request.Context()

	// -------------------- 1) read token --------------------
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest,
			apperror.New(apperror.TokenMalformed, "missing share token"))
		return
	}

	payload, err := utils.ParseShareToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			apperror.New(apperror.TokenInvalid, err.Error()))
		return
	}

	pollID := payload.PollID

	// -------------------- 2) load poll --------------------
	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil || poll.State.IsDeleted {
		c.JSON(http.StatusNotFound,
			apperror.New(apperror.PollNotFound, "poll not found"))
		return
	}

	// -------------------- 3) resolve session --------------------
	var userID string
	var sessionID string

	cookie, err := c.Cookie("session_id")
	if err == nil {
		session, _ := repository.GetSession(cookie)
		if session != nil {
			userID = session.UserID
			sessionID = cookie
		}
	}

	// -------------------- 4) whitelist view restriction --------------------
	if payload.Mode == "whitelist" && userID == "" {
		c.JSON(http.StatusUnauthorized,
			apperror.New(apperror.PollLoginRequired, "login required to view this poll"))
		return
	}

	// -------------------- 5) create visitor session --------------------
	if sessionID == "" {

		anonSession, err := repository.CreateAnonymousSession(ctx, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				apperror.New(apperror.SessionCreateFailed, "failed to create session"))
			return
		}

		sessionID = anonSession.SessionID

		isProd := os.Getenv("APP_ENV") == "prod"

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   int((30 * 24 * time.Hour).Seconds()),
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
	}

	// -------------------- 6) compute viewer state --------------------
	viewer, _ := services.BuildViewerState(ctx, poll, userID, sessionID)

	// -------------------- 7) response --------------------
	c.JSON(http.StatusOK, gin.H{
		"poll": poll,
		"viewer": gin.H{
			"user_id":          userID,
			"session_id":       sessionID,
			"can_vote":         viewer.CanVote,
			"can_view_results": viewer.CanViewResults,
			"voted_option_id":  viewer.VotedOptionID,
			"started":          viewer.Started,
			"ended":            viewer.Ended,
		},
	})
	}
