package ws

import "sync"

type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
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
