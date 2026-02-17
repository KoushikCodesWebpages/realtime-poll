package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(c *gin.Context) {

	pollID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	sess := resolveSession(c, pollID)

	client := newConnection(conn, sess)
	room := getRoom(pollID)

	room.join <- client

	go client.writePump()
}

// ---- PUBLIC EVENT API (used by services) ----

// EmitVote broadcasts a vote delta to the poll room.
// delta must be +1 or -1
func EmitVote(pollID string, optionID string, delta int) {

	r := getRoom(pollID)

	r.broadcast <- internalEvent{
		Type: EventVoteDelta,
		Data: map[string]any{
			"option_id": optionID,
			"delta":     delta,
		},
	}
}

// EmitClosed broadcasts final results when poll closes
func EmitClosed(pollID string, results any) {

	r := getRoom(pollID)

	r.broadcast <- internalEvent{
		Type: EventClosed,
		Data: results,
	}
}