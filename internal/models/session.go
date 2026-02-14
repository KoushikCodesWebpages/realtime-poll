package models

import "time"

type Session struct {
	SessionID string    `bson:"_id"`   // PRIMARY KEY
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}