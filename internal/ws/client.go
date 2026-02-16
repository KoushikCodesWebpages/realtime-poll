package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	room   *Room
	send   chan []byte
	userID string
}

func NewClient(conn *websocket.Conn, room *Room, userID string) *Client {
	return &Client{
		conn:   conn,
		room:   room,
		send:   make(chan []byte, 256),
		userID: userID,
	}
}

func (c *Client) ReadPump(onVote func(userID, optionID, pollID string)) {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg IncomingMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "vote" {
			onVote(c.userID, msg.OptionID, c.room.pollID)
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("ws write error:", err)
			return
		}
	}
}
