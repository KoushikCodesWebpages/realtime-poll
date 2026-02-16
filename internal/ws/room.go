package ws

type Room struct {
	pollID string

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	clients map[*Client]bool
}

func NewRoom(pollID string) *Room {
	return &Room{
		pollID:     pollID,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:
			r.clients[client] = true

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (r *Room) Broadcast(msg []byte) {
	r.broadcast <- msg
}

func (r *Room) Join(c *Client) {
	r.register <- c
}

func (r *Room) Leave(c *Client) {
	r.unregister <- c
}
