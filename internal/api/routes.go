package api

import (
	"github.com/gin-gonic/gin"

	"realtime-poll/internal/middleware"
	"realtime-poll/internal/apperror"
	"realtime-poll/internal/services"
	"realtime-poll/internal/utils"
	"realtime-poll/internal/ws"
)

func RegisterRoutes(
	r *gin.Engine,
	voteService *services.VoteService,
	hub *ws.Hub,
) {

	r.GET("/test-error", func(c *gin.Context) {
		middleware.Fail(c, apperror.Unauthorized())
	})

	// Root info
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, utils.GetRootDoc())
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/b1")
	{
		// ---------------- AUTH ----------------
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.POST("/logout", Logout)
			auth.GET("/session", middleware.RequireAuth(), Session)
		}

		// ---------------- POLL ----------------
		poll := api.Group("/poll", middleware.RequireAuth())
		{
			poll.POST("/create", CreatePoll)
			poll.GET("/mine", GetMyPolls)
			poll.GET("/:poll_id", GetPoll)

			poll.PUT("/:poll_id", PutPoll)
			poll.PATCH("/:poll_id", PatchPoll)
			poll.DELETE("/:poll_id", DeletePoll)

			poll.POST("/:poll_id/share", GenerateShareLink)
		}

		// ---------------- VOTING ----------------
		api.POST("/vote", middleware.OptionalAuth(), CastVote)
		api.GET("/poll/share", ViewSharedPoll)
	}

	// ---------------- WEBSOCKET ----------------
	r.GET("/ws/poll/:pollId", WsPoll(hub, voteService))
}
