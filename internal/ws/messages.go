package ws

type IncomingMessage struct {
	Type     string `json:"type"`
	OptionID string `json:"option_id,omitempty"`
}

type ServerMessage struct {
	Type    string      `json:"type"`
	Results interface{} `json:"results,omitempty"`
	Poll    interface{} `json:"poll,omitempty"`
}