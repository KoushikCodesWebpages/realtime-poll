package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Room string
}
