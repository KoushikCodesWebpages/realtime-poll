package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(c *gin.Context) {
	room := c.Param("id")

	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	client := &Client{Conn: conn, Room: room}

	H.Join(room, client)
}
