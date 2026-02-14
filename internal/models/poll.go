package models

import "time"

type Option struct {
	ID    string `bson:"id" json:"id"`
	Text  string `bson:"text" json:"text"`
	Votes int    `bson:"votes" json:"votes"`
}

type Poll struct {
	ID        string    `bson:"_id" json:"id"`
	Question  string    `bson:"question" json:"question"`
	Options   []Option  `bson:"options" json:"options"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	IsClosed  bool      `bson:"is_closed" json:"is_closed"`
}

