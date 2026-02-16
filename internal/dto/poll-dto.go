package dto

import (
	"time"

	"realtime-poll/internal/models"
)



type VoteUpdate struct {
	Type    string        `json:"type"`
	PollID  string        `json:"poll_id"`
	Results []OptionResult `json:"results"`
}

type OptionResult struct {
	OptionID string `json:"option_id"`
	Votes    int    `json:"votes"`
}


type CreatePollReq struct {
	Question    string   `json:"question" binding:"required,min=5,max=200"`
	Description string   `json:"description,omitempty"`
	Options     []string `json:"options" binding:"required,min=2,max=10"`

	Visibility    string   `json:"visibility"` // public | authenticated | whitelist | link
	AllowedEmails []string `json:"allowed_emails,omitempty"`

	AllowChange bool `json:"allow_change"`
	Anonymous   bool `json:"anonymous"`

	StartAt *time.Time `json:"start_at,omitempty"`
	EndAt   *time.Time `json:"end_at,omitempty"`
}

type PollListResponse struct {
	Data  []models.Poll `json:"data"`
	Next  *string       `json:"next"`
	Prev  *string       `json:"prev"`
	Total int64         `json:"total"`
}

type PollFilter struct {
	Status     string     // active | scheduled | closed | expired
	Search     string
	Visibility string

	HasVotes *bool

	DateFrom *time.Time
	DateTo   *time.Time

	Sort  string // newest | oldest | most_voted
	Limit int64
}