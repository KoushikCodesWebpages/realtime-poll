package ws

import "sync"

type Hub struct {
	mu sync.Mutex

	clients map[*Client]bool

	// one websocket per session
	sessionIndex map[string]*Client

	// poll rooms
	rooms map[string]*Room

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[*Client]bool),
		sessionIndex: make(map[string]*Client),
		rooms:        make(map[string]*Room),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan []byte),
	}
}


func (h *Hub) GetRoom(pollID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[pollID]; ok {
		return room
	}

	room := NewRoom(pollID)
	h.rooms[pollID] = room
	go room.Run()

	return room
}

func (h *Hub) RemoveRoom(pollID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, pollID)
}

func (h *Hub) Run() {
	for {
		select {

		// -------- CONNECT --------
		case client := <-h.register:

			// if same session already connected -> kick old socket
			if old, ok := h.sessionIndex[client.sessionID]; ok {
				old.conn.Close()
				delete(h.clients, old)
			}

			h.clients[client] = true
			h.sessionIndex[client.sessionID] = client

		// -------- DISCONNECT --------
		case client := <-h.unregister:

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

				// remove session mapping only if same instance
				if h.sessionIndex[client.sessionID] == client {
					delete(h.sessionIndex, client.sessionID)
				}

				client.conn.Close()
			}

		// -------- GLOBAL BROADCAST --------
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}