package api

import (
	"os"
	"net/http"
	"time"


	"github.com/gin-gonic/gin"
	"realtime-poll/internal/services"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/constants"
	// "realtime-poll/config"
)

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`

}

type LoginReq struct {
	Identifier string `json:"identifier"` // username OR email
	Password   string `json:"password"`
}
func Session(c *gin.Context) {
	uid := c.GetString(constants.CtxUserID)

	user, err := repository.FindUserByAuthID(uid)
	if err != nil || user == nil {
		c.JSON(401, gin.H{"error": "not authenticated"})
		return
	}

	c.JSON(200, gin.H{
		"id": user.AuthUserID,
		"username": user.Username,
		"email": user.Email,
	})
}



func Register(c *gin.Context) {
	var req AuthReq
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	if err := services.Register(req.Username,req.Email, req.Password); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "registered"})
}

func Login(c *gin.Context) {

	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	result, err := services.Login(req.Identifier, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	secure := os.Getenv("APP_ENV") == "prod"
	maxAge := int(time.Until(result.Session.ExpiresAt).Seconds())

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    result.Session.SessionID,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(200, gin.H{
		"status": "logged_in",
		"user": gin.H{
			"id":       result.User.AuthUserID,
			"username": result.User.Username,
			"email":    result.User.Email,
		},
	})
}


func Logout(c *gin.Context) {

	cookie, err := c.Cookie("session_id")
	if err == nil {
		repository.DeleteSession(cookie)
	}

	secure := os.Getenv("APP_ENV") == "prod"

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite:  http.SameSiteNoneMode,
	})

	c.JSON(200, gin.H{"status": "logged_out"})
}

