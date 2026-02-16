package ws

import (
	"log"
	"sync"
)

/*
Hub manages all active rooms.
It does NOT send messages.
Rooms handle broadcasting themselves.
*/
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

/* ---------------- Room Access ---------------- */

// GetRoom returns existing room or creates one
func (h *Hub) GetRoom(pollID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[pollID]
	if exists {
		return room
	}

	log.Println("WS creating room:", pollID)

	room = NewRoom(pollID)
	h.rooms[pollID] = room
	return room
}

/* ---------------- Cleanup ---------------- */

// RemoveRoom deletes empty rooms
func (h *Hub) RemoveRoom(pollID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.rooms, pollID)
	log.Println("WS removed empty room:", pollID)
}

/* ---------------- Debug ---------------- */

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}
