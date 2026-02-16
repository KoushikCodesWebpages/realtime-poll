package ws
import (
	"context"
	"encoding/json"
	"time"

	"realtime-poll/internal/repository"
	"realtime-poll/internal/services"
)

type Room struct {
	pollID string

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	clients map[*Client]bool

	// NEW
	endTimer *time.Timer
	closed   bool
}


func NewRoom(pollID string) *Room {
	return &Room{
		pollID:     pollID,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		
	}
}
func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:

			r.clients[client] = true

			// start lifecycle timer only once
			if len(r.clients) == 1 {
				go r.startEndTimer()
			}

		r.broadcastPresence()


		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				r.broadcastPresence()
			}

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
func (r *Room) broadcastPresence() {

	payload := map[string]any{
		"type":    "presence",
		"viewers": len(r.clients),
	}

	bytes, _ := json.Marshal(payload)

	for client := range r.clients {
		select {
		case client.send <- bytes:
		default:
			delete(r.clients, client)
			close(client.send)
		}
	}
}


func (r *Room) Broadcast(msg []byte) {
	r.broadcast <- msg
}

func (r *Room) Join(c *Client) {
	r.register <- c
}

func (r *Room) Leave(c *Client) {
	r.unregister <- c
}

func (r *Room) startEndTimer() {

	poll, _ := repository.GetPollByID(context.Background(), r.pollID)
	if poll == nil || poll.Behavior.EndAt == nil {
		return
	}

	duration := time.Until(*poll.Behavior.EndAt)

	if duration <= 0 {
		r.closePoll()
		return
	}

	r.endTimer = time.AfterFunc(duration, func() {
		r.closePoll()
	})
}

func (r *Room) BroadcastJSON(v any) {
	bytes, _ := json.Marshal(v)
	r.Broadcast(bytes)
}


func (r *Room) closePoll() {

	if r.closed {
		return
	}
	r.closed = true

	// mark DB closed
	repository.ClosePoll(context.Background(), r.pollID)

	// fetch results
	results, _ := new(services.VoteService).GetResults(context.Background(), r.pollID)

	payload := map[string]any{
		"type":    "final_results",
		"poll_id": r.pollID,
		"results": results,
	}

	bytes, _ := json.Marshal(payload)

	r.Broadcast(bytes)
}
