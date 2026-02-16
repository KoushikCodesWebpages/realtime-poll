package ws

type IncomingMessage struct {
	Type     string `json:"type"`
	OptionID string `json:"option_id,omitempty"`
}

type VoteUpdate struct {
	Type    string        `json:"type"`
	PollID  string        `json:"poll_id"`
	Results []OptionResult `json:"results"`
}

type OptionResult struct {
	OptionID string `json:"option_id"`
	Votes    int    `json:"votes"`
}
