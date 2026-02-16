package ws

type IncomingMessage struct {
	Type     string `json:"type"`
	OptionID string `json:"option_id,omitempty"`
}

