package ws
import(
	"encoding/json"
)
type Room struct {
	pollID string

	clients map[*Client]bool

	Join      chan *Client
	Leave     chan *Client
	Broadcast chan []byte
}

func NewRoom(pollID string) *Room {
	r := &Room{
		pollID:    pollID,
		clients:   make(map[*Client]bool),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Broadcast: make(chan []byte, 256),
	}

	go r.Run()
	return r
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Join:
			r.clients[client] = true

		case client := <-r.Leave:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				client.Close()
			}

			// auto cleanup empty room
			if len(r.clients) == 0 {
				go GlobalHub.RemoveRoom(r.pollID)
				return
			}


		case msg := <-r.Broadcast:
			for client := range r.clients {
				client.Send(msg)
			}
		}
	}
}


func (r *Room) BroadcastJSON(v any) {
	b, _ := json.Marshal(v)
	r.Broadcast <- b
}