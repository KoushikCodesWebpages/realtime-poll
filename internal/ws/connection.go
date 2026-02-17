package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	ws   *websocket.Conn
	send chan Envelope
	sess *Session

	closeOnce sync.Once
}
func newConnection(ws *websocket.Conn, sess *Session) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan Envelope, 32),
		sess: sess,
	}
}

// READ LOOP — detects disconnects
func (c *Connection) readPump(r *Room) {
	defer func() {
		r.leave <- c
	}()

	for {
		if _, _, err := c.ws.ReadMessage(); err != nil {
			break
		}
	}
}

// WRITE LOOP
func (c *Connection) writePump(r *Room) {
	defer func() {
		r.leave <- c
	}()

	for msg := range c.send {
		if err := c.ws.WriteJSON(msg); err != nil {
			break
		}
	}
}


func (c *Connection) close() {
	c.closeOnce.Do(func() {
		close(c.send)
		c.ws.Close()
	})
}
