package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"realtime-poll/internal/apperror"
)

func Fail(c *gin.Context, err error) {

	if appErr, ok := err.(*apperror.AppError); ok {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    appErr.Code,
			"message": appErr.Message,
			"action":  appErr.Action,
		})
		return
	}

	// unknown error
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":    apperror.INTERNAL_ERROR,
		"message": "Unexpected error",
	})
}
