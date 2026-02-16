package ws

import "encoding/json"

func Write(conn *Client, v any) {
	data, _ := json.Marshal(v)
	conn.Send(data)
}
