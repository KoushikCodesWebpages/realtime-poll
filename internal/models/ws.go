package models

type WSMessage struct {
    Type    string      `json:"type"`
    PollID  string      `json:"poll_id,omitempty"`
    Payload interface{} `json:"payload,omitempty"`
}
