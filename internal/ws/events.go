package ws

type EventType string

const (
	EventPollState EventType = "poll_state"
	EventVoteDelta EventType = "vote_delta"
	EventVoteAck   EventType = "vote_ack"
	EventClosed    EventType = "poll_closed"
	EventError     EventType = "error"
)

type Envelope struct {
	Type EventType `json:"type"`
	Data any       `json:"data,omitempty"`
}

type internalEvent struct {
	Type EventType
	Data any
}
