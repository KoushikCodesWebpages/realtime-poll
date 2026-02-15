package api

import (
	"realtime-poll/internal/ws"
	"realtime-poll/internal/utils"
	"realtime-poll/internal/middleware"

	"github.com/gin-gonic/gin"
)
func RegisterRoutes(r *gin.Engine) {

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
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.POST("/logout", Logout)
		}

		// Poll APIs
	poll := api.Group("/poll", middleware.RequireAuth(),)
	{
		poll.POST("/",CreatePoll) // protected
		poll.GET("/mine", GetMyPolls)

			poll.PUT("/:poll_id", PutPoll)
			poll.PATCH("/:poll_id", PatchPoll)

			poll.DELETE("/:poll_id",DeletePoll)

			poll.POST("/:poll_id/share",GenerateShareLink)
		// poll.POST("/:id/vote", Vote) // public
	}

		// Websocket	
		api.GET("/ws/:id", ws.HandleWS)

		r.GET("/b1/poll/share", ViewSharedPoll) 
	}

}
