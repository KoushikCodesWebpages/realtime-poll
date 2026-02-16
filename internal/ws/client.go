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
		c.room.Leave(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			break
		}

		log.Println("WS RECEIVED:", string(message))

		var msg IncomingMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("invalid ws json")
			continue
		}

		if msg.Type == "vote" && msg.OptionID != "" {
			onVote(c.userID, msg.OptionID, c.room.pollID)
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("WS write error:", err)
			return
		}
	}
}

func (c *Client) Send(msg []byte) {
	select {
	case c.send <- msg:
	default:
	}
}
