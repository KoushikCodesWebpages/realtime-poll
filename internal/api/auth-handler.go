package api

import (

	"github.com/gin-gonic/gin"
	"realtime-poll/internal/services"
	"realtime-poll/internal/repository"
)

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req AuthReq
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	if err := services.Register(req.Username, req.Password); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "registered"})
}

func Login(c *gin.Context) {
	var req AuthReq
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	session, err := services.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("session_id", session.ID, 86400, "/", "", false, true)

	c.JSON(200, gin.H{"status": "logged_in"})
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err == nil {
		repository.DeleteSession(cookie)
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{"status": "logged_out"})
}
