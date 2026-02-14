package models

import "time"

type User struct {
	ID           string        `bson:"_id,omitempty" json:"id"`
	Username     string        `bson:"username" json:"username"`
	Email        string		   `bson:"email"  json:"email"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`

	Profile      Profile       `bson:"profile,omitempty" json:"profile,omitempty"`
	Security     Security      `bson:"security,omitempty" json:"security,omitempty"`
	Preferences  Preferences   `bson:"preferences,omitempty" json:"preferences,omitempty"`
	Stats        Stats         `bson:"stats,omitempty" json:"stats,omitempty"`
	Social       Social        `bson:"social,omitempty" json:"social,omitempty"`
	Metadata     Metadata      `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type Profile struct {
	DisplayName   string    `bson:"display_name,omitempty"`
	Bio           string    `bson:"bio,omitempty"`
	AvatarURL     string    `bson:"avatar_url,omitempty"`
	BannerURL     string    `bson:"banner_url,omitempty"`
	Website       string    `bson:"website,omitempty"`
	Location      string    `bson:"location,omitempty"`
	Company       string    `bson:"company,omitempty"`
	JobTitle      string    `bson:"job_title,omitempty"`
	Pronouns      string    `bson:"pronouns,omitempty"`
	Birthdate     time.Time `bson:"birthdate,omitempty"`
	Timezone      string    `bson:"timezone,omitempty"`
	Language      string    `bson:"language,omitempty"`
	Theme         string    `bson:"theme,omitempty"`
}

type Security struct {
	LastLoginAt        time.Time `bson:"last_login_at,omitempty"`
	LastPasswordChange time.Time `bson:"last_password_change,omitempty"`
	FailedLoginCount   int       `bson:"failed_login_count,omitempty"`
	LockedUntil        time.Time `bson:"locked_until,omitempty"`
	TwoFactorEnabled   bool      `bson:"two_factor_enabled,omitempty"`
	RecoveryEmail      string    `bson:"recovery_email,omitempty"`
	RecoveryPhone      string    `bson:"recovery_phone,omitempty"`
}

type Preferences struct {
	EmailNotifications bool   `bson:"email_notifications,omitempty"`
	PushNotifications  bool   `bson:"push_notifications,omitempty"`
	MarketingEmails    bool   `bson:"marketing_emails,omitempty"`
	DefaultPollPrivacy string `bson:"default_poll_privacy,omitempty"`
	AutoClosePolls     bool   `bson:"auto_close_polls,omitempty"`
	CompactMode        bool   `bson:"compact_mode,omitempty"`
	DarkMode           bool   `bson:"dark_mode,omitempty"`
	ShowOnlineStatus   bool   `bson:"show_online_status,omitempty"`
}

type Stats struct {
	PollsCreated    int `bson:"polls_created,omitempty"`
	VotesCast       int `bson:"votes_cast,omitempty"`
	TotalViews      int `bson:"total_views,omitempty"`
	ActivePolls     int `bson:"active_polls,omitempty"`
	ClosedPolls     int `bson:"closed_polls,omitempty"`
	ReportsReceived int `bson:"reports_received,omitempty"`
}

type Social struct {
	Twitter   string `bson:"twitter,omitempty"`
	Github    string `bson:"github,omitempty"`
	LinkedIn  string `bson:"linkedin,omitempty"`
	YouTube   string `bson:"youtube,omitempty"`
	Twitch    string `bson:"twitch,omitempty"`
	Discord   string `bson:"discord,omitempty"`
	Telegram  string `bson:"telegram,omitempty"`
	Instagram string `bson:"instagram,omitempty"`
}

type Metadata struct {
	IsAdmin        bool      `bson:"is_admin,omitempty"`
	IsBanned       bool      `bson:"is_banned,omitempty"`
	BannedReason   string    `bson:"banned_reason,omitempty"`
	BannedUntil    time.Time `bson:"banned_until,omitempty"`
	AccountTier    string    `bson:"account_tier,omitempty"`
	SignupIP       string    `bson:"signup_ip,omitempty"`
	LastActiveIP   string    `bson:"last_active_ip,omitempty"`
	SignupSource   string    `bson:"signup_source,omitempty"`
	FeatureFlags   []string  `bson:"feature_flags,omitempty"`
	InternalNotes  string    `bson:"internal_notes,omitempty"`
}
