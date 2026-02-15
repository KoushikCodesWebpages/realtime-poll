package models

import "time"

type Option struct {
	OptionID    string `bson:"option_id" json:"option_id"`
	Text  string `bson:"text" json:"text"`
	Votes int    `bson:"votes" json:"votes"`
}

type Poll struct {
	PollID  string `bson:"poll_id" json:"poll_id"`
	OwnerID string `bson:"owner_id" json:"owner_id"`

	Content      ContentSettings      `bson:"content" json:"content"`
	Access       AccessSettings       `bson:"access" json:"access"`
	Vote         VoteSettings         `bson:"vote" json:"vote"`
	Distribution DistributionSettings `bson:"distribution" json:"distribution"`
	Behavior     BehaviorSettings     `bson:"behavior" json:"behavior"`
	Analytics    AnalyticsSettings    `bson:"analytics" json:"analytics"`

	State PollState `bson:"state" json:"state"`
	Meta  Meta      `bson:"meta" json:"meta"`
}

type ContentSettings struct {
	Question    string   `bson:"question" json:"question"`
	Description string   `bson:"description,omitempty" json:"description,omitempty"`
	Options     []Option `bson:"options" json:"options"`

	Images      []string `bson:"images,omitempty" json:"images,omitempty"`
	AllowCustom bool     `bson:"allow_custom_option" json:"allow_custom_option"`
	Randomize   bool     `bson:"randomize_options" json:"randomize_options"`
}

type AccessSettings struct {
	Visibility    string   `bson:"visibility" json:"visibility"` // public | authenticated | whitelist | link
	AllowedUsers  []string `bson:"allowed_users,omitempty" json:"allowed_users,omitempty"`
	AllowedEmails []string `bson:"allowed_emails,omitempty" json:"allowed_emails,omitempty"`
	BlockedUsers  []string `bson:"blocked_users,omitempty" json:"blocked_users,omitempty"`

	RequireLogin bool `bson:"require_login" json:"require_login"`
}
type VoteSettings struct {
	MaxVotesPerUser int  `bson:"max_votes_per_user" json:"max_votes_per_user"`
	AllowChangeVote bool `bson:"allow_change_vote" json:"allow_change_vote"`
	AnonymousVote   bool `bson:"anonymous_vote" json:"anonymous_vote"`
	HideResults     bool `bson:"hide_results_until_end" json:"hide_results_until_end"`
	ShowVoters      bool `bson:"show_voters" json:"show_voters"`

	UniqueIP      bool `bson:"unique_ip" json:"unique_ip"`
	UniqueSession bool `bson:"unique_session" json:"unique_session"`
}


type DistributionSettings struct {
	ShareID     string `bson:"share_id,omitempty" json:"share_id,omitempty"`
	QRCode      bool   `bson:"qr_code" json:"qr_code"`
	Embeddable  bool   `bson:"embeddable" json:"embeddable"`
	OneTimeLink bool   `bson:"one_time_link" json:"one_time_link"`

	AllowRepost bool `bson:"allow_repost" json:"allow_repost"`
	AllowExport bool `bson:"allow_export_results" json:"allow_export_results"`
}


type AnalyticsSettings struct {
	TrackViews    bool `bson:"track_views" json:"track_views"`
	TrackVoters   bool `bson:"track_voters" json:"track_voters"`
	TrackLocation bool `bson:"track_location" json:"track_location"`
	TrackDevice   bool `bson:"track_device" json:"track_device"`

	FraudDetection bool `bson:"fraud_detection" json:"fraud_detection"`
}


type BehaviorSettings struct {
	StartAt *time.Time `bson:"start_at,omitempty" json:"start_at,omitempty"`
	EndAt   *time.Time `bson:"end_at,omitempty" json:"end_at,omitempty"`
	AutoClose bool     `bson:"auto_close" json:"auto_close"`

	ShowLiveResults bool `bson:"show_live_results" json:"show_live_results"`
	NotifyOwner     bool `bson:"notify_owner_on_vote" json:"notify_owner_on_vote"`
}


type PollState struct {
	IsClosed bool `bson:"is_closed" json:"is_closed"`
	IsLocked bool `bson:"is_locked" json:"is_locked"`
	Version  int  `bson:"version" json:"version"`
}

type Meta struct {
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	LastVote  *time.Time `bson:"last_vote,omitempty" json:"last_vote,omitempty"`

	TotalVotes int `bson:"total_votes" json:"total_votes"`
	TotalViews int `bson:"total_views" json:"total_views"`
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"` // derived from Behavior.EndAt
}
