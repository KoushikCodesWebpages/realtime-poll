package ws

import (
	"sync"
)

type Room struct {
	PollID    string
	clients   map[*Connection]struct{}
	join      chan *Connection
	leave     chan *Connection
	broadcast chan internalEvent
	closed    bool
	mu        sync.RWMutex
}

func newRoom(pollID string) *Room {
	r := &Room{
		PollID:    pollID,
		clients:   make(map[*Connection]struct{}),
		join:      make(chan *Connection),
		leave:     make(chan *Connection),
		broadcast: make(chan internalEvent, 32),
	}

	go r.run()
	go startExpiryTimer(r)

	return r
}

func (r *Room) run() {
	for {
		select {

		case c := <-r.join:
			r.clients[c] = struct{}{}

		case c := <-r.leave:
			delete(r.clients, c)
			c.close()

		case ev := <-r.broadcast:
			r.emit(ev)
		}
	}
}
