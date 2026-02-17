package api

import (
	// "encoding/json"
	
	// "context"
	// "log"
	// "time"
	// "realtime-poll/internal/repository"
	// "realtime-poll/internal/services"
	
	// "realtime-poll/internal/utils"
	// "realtime-poll/internal/apperror"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"realtime-poll/internal/ws"
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}


func ServePollWS(c *gin.Context) {
	ws.ServeWS(c)
}
