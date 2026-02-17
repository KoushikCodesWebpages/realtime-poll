package dto

import (
	"time"

	"realtime-poll/internal/models"
)

type EditPollReq struct {
	Content  *EditContent  `json:"content"`
	Access   *EditAccess   `json:"access"`
	Vote     *EditVote     `json:"vote"`
	Behavior *EditBehavior `json:"behavior"`
}

type EditContent struct {
	Question    *string       `json:"question,omitempty"`
	Description *string       `json:"description,omitempty"`
	Options     []EditOption  `json:"options,omitempty"`

	AllowCustomOption *bool `json:"allow_custom_option,omitempty"`
	RandomizeOptions  *bool `json:"randomize_options,omitempty"`
}

type EditOption struct {
	OptionID string `json:"option_id"`
	Text     string `json:"text"`
}

type EditAccess struct {
	Visibility    *string  `json:"visibility,omitempty"`   // public | authenticated | whitelist | link
	AllowedEmails []string `json:"allowed_emails,omitempty"`
}

type EditVote struct {
	MaxVotesPerUser   *int  `json:"max_votes_per_user,omitempty"`
	AllowChangeVote   *bool `json:"allow_change_vote,omitempty"`
	AnonymousVote     *bool `json:"anonymous_vote,omitempty"`
	HideResultsUntilEnd *bool `json:"hide_results_until_end,omitempty"`
	ShowVoters        *bool `json:"show_voters,omitempty"`
	UniqueIP          *bool `json:"unique_ip,omitempty"`
	UniqueSession     *bool `json:"unique_session,omitempty"`
}

type EditBehavior struct {
	StartAt         *time.Time `json:"start_at,omitempty"`
	EndAt           *time.Time `json:"end_at,omitempty"`
	AutoClose       *bool      `json:"auto_close,omitempty"`
	ShowLiveResults *bool      `json:"show_live_results,omitempty"`
	NotifyOwner     *bool      `json:"notify_owner_on_vote,omitempty"`
}


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

	// ================= CONTENT =================
	Question    string   `json:"question" binding:"required,min=5,max=200"`
	Description string   `json:"description,omitempty"`
	Options     []string `json:"options" binding:"required,min=2,max=10"`

	Images           []string `json:"images,omitempty"`
	AllowCustomOption bool    `json:"allow_custom_option"`
	RandomizeOptions  bool    `json:"randomize_options"`

	// ================= ACCESS =================
	Visibility    string   `json:"visibility" binding:"oneof=public authenticated whitelist link"`
	AllowedEmails []string `json:"allowed_emails,omitempty"`

	// ================= VOTING =================
	MaxVotesPerUser int  `json:"max_votes_per_user"`
	AllowChange     bool `json:"allow_change"`
	Anonymous       bool `json:"anonymous"`

	HideResultsUntilEnd bool `json:"hide_results_until_end"`
	ShowVoters          bool `json:"show_voters"`

	UniqueIP      bool `json:"unique_ip"`
	UniqueSession bool `json:"unique_session"`

	// ================= BEHAVIOR =================
	StartAt *time.Time `json:"start_at,omitempty"`
	EndAt   *time.Time `json:"end_at,omitempty"`

	AutoClose       bool `json:"auto_close"`
	ShowLiveResults bool `json:"show_live_results"`
	NotifyOwnerOnVote bool `json:"notify_owner_on_vote"`

	// ================= DISTRIBUTION =================
	ShareType string `json:"share_type" binding:"omitempty,oneof=none link qr embed"`

	// ================= ANALYTICS (optional future) =================
	TrackLocation bool `json:"track_location"`
	TrackDevice   bool `json:"track_device"`
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