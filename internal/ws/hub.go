package ws

type Hub struct {
	rooms map[string]map[*Client]bool
}

var H = Hub{
	rooms: make(map[string]map[*Client]bool),
}

func (h *Hub) Join(room string, c *Client) {
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][c] = true
}

func (h *Hub) Broadcast(room string, msg []byte) {
	for c := range h.rooms[room] {
		c.Conn.WriteMessage(1, msg)
	}
}
