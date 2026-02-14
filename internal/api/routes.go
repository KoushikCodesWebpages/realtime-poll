package api

import (
	"realtime-poll/internal/ws"
	"realtime-poll/internal/utils"

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
		poll := api.Group("/poll")
		{
			poll.POST("", CreatePoll)
			poll.POST("/:id/vote", Vote)
		}

		// Websocket
		api.GET("/ws/:id", ws.HandleWS)
	}

}
