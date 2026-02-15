package dto

import "time"

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