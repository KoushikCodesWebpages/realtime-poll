package ws

import (
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws   *websocket.Conn
	send chan Envelope
	sess *Session
}

func newConnection(ws *websocket.Conn, sess *Session) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan Envelope, 32),
		sess: sess,
	}
}

func (c *Connection) writePump() {
	defer c.ws.Close()

	for msg := range c.send {
		_ = c.ws.WriteJSON(msg)
	}
}

func (c *Connection) close() {
	close(c.send)
	c.ws.Close()
}
