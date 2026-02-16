package api

import (
	"encoding/json"
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"realtime-poll/internal/repository"
	"realtime-poll/internal/services"
	"realtime-poll/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsPoll(hub *ws.Hub, voteService *services.VoteService) gin.HandlerFunc {
	return func(c *gin.Context) {

		pollID := c.Param("pollId")

		// -------------------- AUTH --------------------
		sessionCookie, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(401, gin.H{"error": "not logged in"})
			return
		}

		session, err := repository.GetSession(sessionCookie)
		if err != nil || session == nil {
			c.JSON(401, gin.H{"error": "invalid session"})
			return
		}

		userID := session.UserID
		ip := c.ClientIP()
		

		// -------------------- UPGRADE --------------------
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// -------------------- ROOM JOIN --------------------
		room := hub.GetRoom(pollID)
		client := ws.NewClient(conn, room, userID)

		room.Join(client)

		go client.WritePump()

		// -------------------- READ LOOP --------------------
		go client.ReadPump(func(userID, optionID, pollID string) {

		// NEVER use gin request context after upgrade
		ctx := context.Background()

		err := voteService.CastVoteRealtime(
			ctx,
			pollID,
			optionID,
			userID,
			sessionCookie,
			ip,
		)
		if err != nil {
			println("VOTE ERROR:", err.Error())
			return
		}

		results, err := voteService.GetResults(ctx, pollID)
		if err != nil {
			println("RESULT ERROR:", err.Error())
			return
		}

		payload := ws.VoteUpdate{
			Type:    "vote_update",
			PollID:  pollID,
			Results: results,
		}

		bytes, _ := json.Marshal(payload)
		room.Broadcast(bytes)
	})
	}
}
