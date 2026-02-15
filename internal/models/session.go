package models

import "time"

type Session struct {
	SessionID string    `bson:"_id"`   // PRIMARY KEY
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

type ShareTokenClaims struct {
	PollID string `json:"poll"`
	Type   string `json:"type"` // view
	Uses   int    `json:"uses"` // allowed uses
	Exp    int64  `json:"exp"`
}