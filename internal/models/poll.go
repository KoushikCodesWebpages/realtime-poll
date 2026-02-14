package models

import "time"

type Option struct {
	OptionID    string `bson:"option_id" json:"option_id"`
	Text  string `bson:"text" json:"text"`
	Votes int    `bson:"votes" json:"votes"`
}

type Poll struct {
	AuthUserID    string    `bson:"auth_user_id,omitempty" json:"auth_user_id"`
	PollID        string    `bson:"poll_id" json:"poll_id"`
	Question  string    `bson:"question" json:"question"`
	Options   []Option  `bson:"options" json:"options"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	IsClosed  bool      `bson:"is_closed" json:"is_closed"`
}

type CreatePollReq struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

