package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/*
Client represents one connected websocket user
*/
type Client struct {
	conn *websocket.Conn
	room *Room

	send chan []byte

	userID string

	closeOnce sync.Once
}

/* ---------------- Constructor ---------------- */

func NewClient(conn *websocket.Conn, room *Room, userID string) *Client {
	return &Client{
		conn:   conn,
		room:   room,
		send:   make(chan []byte, 256),
		userID: userID,
	}
}

/* ---------------- Close ---------------- */

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		log.Println("WS client closing")
		close(c.send)
		c.conn.Close()
	})
}

/* ---------------- Read Loop ---------------- */

func (c *Client) ReadPump(onVote func(userID, optionID, pollID string)) {

	defer func() {
		c.room.Leave <- c
	}()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			return
		}

		var msg IncomingMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			SendError(c, err)
			continue
		}

		switch msg.Type {

		case "ping":
			Write(c, WSResponse{
				Type:    "pong",
				Success: true,
			})

		case "vote":
			if msg.OptionID == "" {
				SendError(c, ErrInvalidVotePayload)
				continue
			}

			onVote(c.userID, msg.OptionID, c.room.pollID)

		default:
			SendError(c, ErrUnknownMessageType)
		}
	}
}

/* ---------------- Write Loop ---------------- */

func (c *Client) WritePump() {

	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		c.room.Leave <- c
	}()

	for {
		select {

		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("WS write error:", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/* ---------------- Send Helpers ---------------- */

func (c *Client) Send(msg []byte) {

	defer func() {
		if recover() != nil {
			// channel already closed
			log.Println("WS send on closed client ignored")
		}
	}()

	select {
	case c.send <- msg:
	default:
		log.Println("WS slow client dropped")
		c.room.Leave <- c
	}
}


func (c *Client) SendJSON(v any) {
	b, _ := json.Marshal(v)
	c.Send(b)
}
