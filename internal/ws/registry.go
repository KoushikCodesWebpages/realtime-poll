package ws

import "sync"

type registry struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

var hub = &registry{
	rooms: make(map[string]*Room),
}

func getRoom(pollID string) *Room {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if r, ok := hub.rooms[pollID]; ok {
		return r
	}

	r := newRoom(pollID)
	hub.rooms[pollID] = r
	return r
}
