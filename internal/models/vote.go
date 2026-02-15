package models

import "time"
type VoteRecord struct {
	VoteID   string `bson:"vote_id"`
	PollID   string `bson:"poll_id"`

	UserID   string `bson:"user_id,omitempty"`
	SessionID string `bson:"session_id,omitempty"`
	IPAddress string `bson:"ip_address,omitempty"`

	OptionID string `bson:"option_id"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
