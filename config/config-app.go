package config

import (
	"encoding/json"
	"log"
	"os"
)

var AppRules AppConfig

type AppConfig struct {
	Poll      PollConfig      `json:"poll"`
	Voting    VotingConfig    `json:"voting"`
	AntiAbuse AntiAbuseConfig `json:"anti_abuse"`
	Websocket WSConfig        `json:"websocket"`
	RateLimit RateLimitConfig `json:"rate_limit"`
	Storage   StorageConfig   `json:"storage"`
	Logging   LoggingConfig   `json:"logging"`
}

type PollConfig struct {
	MaxOptionsPerPoll int  `json:"max_options_per_poll"`
	AllowVoteChange   bool `json:"allow_vote_change"`
}

type VotingConfig struct {
	IPVoteLimit          int `json:"ip_vote_limit"`
	VoteCooldownSeconds  int `json:"vote_cooldown_seconds"`
	OneVotePerFingerprint bool `json:"one_vote_per_fingerprint"`
}

type AntiAbuseConfig struct {
	RapidVoteThreshold     int `json:"rapid_vote_threshold"`
	RapidVoteWindowSeconds int `json:"rapid_vote_window_seconds"`
	TemporaryBlockMinutes  int `json:"temporary_block_minutes"`
}

type WSConfig struct {
	MaxConnectionsPerRoom int `json:"max_connections_per_room"`
}

type RateLimitConfig struct {
	CreatePollPerMinute int `json:"create_poll_per_minute"`
	VotePerMinutePerIP  int `json:"vote_per_minute_per_ip"`
}

type StorageConfig struct {
	StoreIPHashOnly bool `json:"store_ip_hash_only"`
}

type LoggingConfig struct {
	LogVotes bool `json:"log_votes"`
}

func LoadAppConfig() {
	file, err := os.ReadFile("configs/app.json")
	if err != nil {
		log.Fatal("cannot load app.json")
	}

	err = json.Unmarshal(file, &AppRules)
	if err != nil {
		log.Fatal("invalid app.json")
	}

	log.Println("App config loaded")
}
