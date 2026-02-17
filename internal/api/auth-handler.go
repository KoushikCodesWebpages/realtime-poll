package api

import (
	"os"
	"net/http"
	"time"


	"github.com/gin-gonic/gin"
	"realtime-poll/internal/services"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
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

func GetWSToken(c *gin.Context) {
	userID := c.GetString(constants.CtxUserID)

	token, _ := utils.GenerateWSToken(userID)

	c.JSON(200, gin.H{
		"ws_token": token,
	})
}



func Session(c *gin.Context) {
	uid := c.GetString(constants.CtxUserID)

	if uid == "" {
		c.JSON(200, gin.H{"user": nil})
		return
	}

	user, err := repository.FindUserByAuthID(uid)
	if err != nil || user == nil {
		c.JSON(200, gin.H{"user": nil})
		return
	}
	
	c.JSON(200, gin.H{
		"user": gin.H{
			"id":       user.AuthUserID,
			"username": user.Username,
			"email":    user.Email,
		},
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

	maxAge := int(time.Until(result.Session.ExpiresAt).Seconds())

	isProd := os.Getenv("APP_ENV") == "prod"
	maxAge = int(time.Until(result.Session.ExpiresAt).Seconds())
	expires := time.Now().UTC().Add(time.Duration(maxAge) * time.Second)


	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    result.Session.SessionID,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expires,   
		HttpOnly: true,
	}

	if isProd {
		cookie.Domain = ".clqit.in"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		// localhost compatible
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, cookie)

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

	// ---------------- remove session from DB ----------------
	cookieValue, err := c.Cookie("session_id")
	if err == nil && cookieValue != "" {
		repository.DeleteSession(cookieValue)
	}

	// ---------------- build deletion cookie ----------------
	isProd := os.Getenv("APP_ENV") == "prod"

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // delete cookie
		HttpOnly: true,
	}

	if isProd {
		// production domain cookie
		cookie.Domain = ".clqit.in"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		// localhost cookie
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(200, gin.H{
		"status": "logged_out",
	})
}

