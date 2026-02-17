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

// func handleRealtimeVote(
// 	ctx context.Context,
// 	room *ws.Room,
// 	pollID string,
// 	userID string,
// 	optionID string,
// ) {
// 	voteService := services.VoteService{}

// 	version, err := voteService.CastVoteRealtime(
// 		ctx,
// 		pollID,
// 		optionID,
// 		userID,
// 		"", // session handled in HTTP only
// 		"",
// 	)

// 	if err != nil {
// 		room.BroadcastJSON(gin.H{
// 			"type": "vote_rejected",
// 			"code": err.(*apperror.AppError).Code,
// 			"message": err.(*apperror.AppError).Message,
// 		})
// 		return
// 	}

// 	// fetch updated counts
// 	results, _ := voteService.GetResults(ctx, pollID)

// 	for _, r := range results {
// 		room.BroadcastJSON(gin.H{
// 			"type":       "vote_update",
// 			"option_id":  r.OptionID,
// 			"votes":      r.Votes,
// 			"total_votes": version,
// 		})
// 	}
// }

// func WsPoll(hub *ws.Hub, voteService *services.VoteService) gin.HandlerFunc {

// 	return func(c *gin.Context) {

// 		pollID := c.Param("pollId")

// 		/* ---------------- AUTH TOKEN ---------------- */

// 		token := c.Query("token")
// 		if token == "" {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "missing ws token"})
// 			return
// 		}

// 		userID, err := utils.ParseWSToken(token)
// 		if err != nil {
// 			log.Println("WS TOKEN ERROR:", err)
// 			c.AbortWithStatusJSON(401, gin.H{"error": "invalid ws token"})
// 			return
// 		}

// 		sessionID := "ws:" + userID
// 		ip := c.ClientIP()

// 		/* ---------------- UPGRADE ---------------- */

// 		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 		if err != nil {
// 			return
// 		}

// 		/* ---------------- CLIENT + ROOM ---------------- */

// 		room := hub.GetRoom(pollID)
// 		client := ws.NewClient(conn, room, userID)

// 		room.Join <- client

// 		/* ---------------- START PUMPS ---------------- */

// 		go client.WritePump()

// 		go client.ReadPump(func(userID, optionID, pollID string) {

// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()

// 			/* -------- load poll -------- */

// 			poll, err := repository.GetPollByID(ctx, pollID)
// 			if err != nil || poll == nil {
// 				sendReject(client, "POLL_NOT_FOUND", "poll not found")
// 				return
// 			}

// 			/* -------- lifecycle -------- */

// 			if poll.State.IsClosed ||
// 				(poll.Behavior.EndAt != nil && time.Now().After(*poll.Behavior.EndAt)) {

// 				sendReject(client, "POLL_ENDED", "poll has ended")
// 				return
// 			}

// 			/* -------- realtime capability -------- */

// 			if !poll.Behavior.ShowLiveResults {
// 				sendReject(client, "NOT_REALTIME", "this poll does not support live voting")
// 				return
// 			}

// 			/* -------- CAST VOTE -------- */

// 			version, err := voteService.CastVoteRealtime(
// 				ctx,
// 				pollID,
// 				optionID,
// 				userID,
// 				sessionID,
// 				ip,
// 			)

// 			if err != nil {
// 				if appErr, ok := err.(*apperror.AppError); ok {
// 					sendReject(client, string(appErr.Code), appErr.Message)
// 				} else {
// 					sendReject(client, "VOTE_FAILED", "vote failed")
// 				}
// 				return
// 			}

// 			now := time.Now()
// 			pollEnded := poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt)

// 			/* -------- hidden ballot -------- */

// 			if poll.Vote.HideResults && !pollEnded {
// 				client.SendJSON(map[string]any{
// 					"type":    "vote_ack",
// 					"version": version,
// 				})
// 				return
// 			}

// 			/* -------- FETCH RESULTS -------- */

// 			results, err := voteService.GetResults(ctx, pollID)
// 			if err != nil {
// 				sendReject(client, "RESULT_ERROR", "failed to fetch results")
// 				return
// 			}

// 			/* -------- BROADCAST -------- */

// 			room.BroadcastJSON(map[string]any{
// 				"type":    "vote_update",
// 				"poll_id": pollID,
// 				"version": version,
// 				"results": results,
// 			})
// 		})

// 		/* ---------------- INITIAL STATE ---------------- */

// 		go sendInitialState(client, voteService, pollID, userID, sessionID)
// 	}
// }

// func sendReject(client *ws.Client, code string, message string) {
// 	payload := map[string]any{
// 		"type":    "vote_rejected",
// 		"code":    code,
// 		"message": message,
// 	}
// 	client.SendJSON(payload)
// }


// func sendInitialState(client *ws.Client, voteService *services.VoteService, pollID, userID, sessionID string) {

// 	ctx := context.Background()

// 	poll, err := repository.GetPollByID(ctx, pollID)
// 	if err != nil || poll == nil {
// 		return
// 	}

// 	// lifecycle
// 	now := time.Now()
// 	started := poll.Behavior.StartAt == nil || now.After(*poll.Behavior.StartAt)
// 	ended := poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt)

// 	// viewer vote
// 	vote, _ := repository.GetVoteForViewer(ctx, pollID, userID, sessionID)

// 	var myVote any = nil
// 	if vote != nil {
// 		myVote = vote.OptionID
// 	}

// 	// results
// 	results, _ := voteService.GetResults(ctx, pollID)

// 	payload := map[string]any{
// 		"type":    "init_state",
// 		"my_vote": myVote,
// 		"results": results,
// 		"started": started,
// 		"ended":   ended,
// 	}

// 	client.SendJSON(payload)
// }
