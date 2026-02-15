package models

import
(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
) 

type Session struct {
	SessionID string    `bson:"_id"`   // PRIMARY KEY
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

type ShareTokenClaims struct {
	PollID string `json:"poll"`
	Mode   string `json:"mode"` // infinite | timed
	Type   string `json:"type"` // view
	Uses   int    `json:"uses"` // allowed uses
	Exp    int64  `json:"exp"`
}

type Vote struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	PollID   string `bson:"poll_id"`
	OptionID string `bson:"option_id"`

	UserID    string `bson:"user_id,omitempty"`
	SessionID string `bson:"session_id,omitempty"`
	IP        string `bson:"ip,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
}