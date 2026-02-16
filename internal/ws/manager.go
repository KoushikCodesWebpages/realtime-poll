package ws

import "sync"

var rooms = struct {
	sync.RWMutex
	m map[string]*Room
}{
	m: make(map[string]*Room),
}

// GetRoom returns existing room or creates one
func GetRoom(pollID string) *Room {
	rooms.Lock()
	defer rooms.Unlock()

	r, ok := rooms.m[pollID]
	if !ok {
		r = NewRoom(pollID)
		rooms.m[pollID] = r
	}
	return r
}
