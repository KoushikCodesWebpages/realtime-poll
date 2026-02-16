package api

import (
	// "encoding/json"
	"net/http"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"realtime-poll/internal/repository"
	"realtime-poll/internal/services"
	"realtime-poll/internal/ws"
	"realtime-poll/internal/dto"
	"realtime-poll/internal/apperror"
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
		go sendInitialState(client, voteService, pollID, userID, sessionCookie)


		// -------------------- READ LOOP --------------------
		go client.ReadPump(func(userID, optionID, pollID string) {

			ctx := context.Background()

			// load poll
			poll, err := repository.GetPollByID(ctx, pollID)
			if err != nil || poll == nil {
				sendReject(client, "POLL_NOT_FOUND", "poll not found")
				return
			}

			// lifecycle FIRST
			if poll.State.IsClosed ||
			(poll.Behavior.EndAt != nil && time.Now().After(*poll.Behavior.EndAt)) {

				sendReject(client, "POLL_ENDED", "poll has ended")
				return
			}

			// realtime capability AFTER lifecycle
			if !poll.Behavior.ShowLiveResults {
				sendReject(client, "NOT_REALTIME", "this poll does not support live voting")
				return
			}
			// cast vote
			err = voteService.CastVoteRealtime(
				ctx,
				pollID,
				optionID,
				userID,
				sessionCookie,
				ip,
			)

			if err != nil {

				// structured apperror
				if appErr, ok := err.(*apperror.AppError); ok {
					sendReject(client, string(appErr.Code), appErr.Message)
				} else {
					sendReject(client, "VOTE_FAILED", err.Error())
				}
				return
			}

			now := time.Now()
			pollEnded := poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt)

			// hidden ballot
			if poll.Vote.HideResults && !pollEnded {

				client.SendJSON(map[string]any{
					"type": "vote_ack",
				})
				return
			}

			// results
			results, err := voteService.GetResults(ctx, pollID)
			if err != nil {
				sendReject(client, "RESULT_ERROR", "failed to fetch results")
				return
			}

			room.BroadcastJSON(dto.VoteUpdate{
				Type:    "vote_update",
				PollID:  pollID,
				Results: results,
			})
		})
	}
}

func sendReject(client *ws.Client, code string, message string) {
	payload := map[string]any{
		"type":    "vote_rejected",
		"code":    code,
		"message": message,
	}
	client.SendJSON(payload)
}


func sendInitialState(client *ws.Client, voteService *services.VoteService, pollID, userID, sessionID string) {

	ctx := context.Background()

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil {
		return
	}

	// lifecycle
	now := time.Now()
	started := poll.Behavior.StartAt == nil || now.After(*poll.Behavior.StartAt)
	ended := poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt)

	// viewer vote
	vote, _ := repository.GetVoteForViewer(ctx, pollID, userID, sessionID)

	var myVote any = nil
	if vote != nil {
		myVote = vote.OptionID
	}

	// results
	results, _ := voteService.GetResults(ctx, pollID)

	payload := map[string]any{
		"type":    "init_state",
		"my_vote": myVote,
		"results": results,
		"started": started,
		"ended":   ended,
	}

	client.SendJSON(payload)
}
