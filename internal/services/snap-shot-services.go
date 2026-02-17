package services

import (
	"context"

	"realtime-poll/internal/apperror"
	"realtime-poll/internal/dto"
	"realtime-poll/internal/repository"

	"time"

	"realtime-poll/internal/ws"
)

type SnapshotService struct{}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{}
}

type PollSnapshot struct {
	PollID     string             `json:"poll_id"`
	Options    []dto.OptionResult `json:"options"`
	TotalVotes int                `json:"total_votes"`
	Version    int64              `json:"version"`
}

func (s *SnapshotService) GetSnapshot(ctx context.Context, pollID string) (*PollSnapshot, error) {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, apperror.Internal()
	}
	if poll == nil {
		return nil, apperror.PollNotFound()
	}

	total := 0
	results := make([]dto.OptionResult, 0, len(poll.Content.Options))

	for _, opt := range poll.Content.Options {
		results = append(results, dto.OptionResult{
			OptionID: opt.OptionID,
			Votes:    opt.Votes,
		})
		total += opt.Votes
	}

	return &PollSnapshot{
		PollID:     pollID,
		Options:    results,
		TotalVotes: total,
		Version:    poll.State.Version,
	}, nil

}


func BuildPollSnapshot(sess *ws.Session) (any, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll, err := repository.GetPollByID(ctx, sess.PollID)
	if err != nil || poll == nil {
		return nil, err
	}

	// ----- owner override -----
	isOwner := false
	if sess.UserID != nil && *sess.UserID == poll.OwnerID {
		isOwner = true
	}

	// determine if votes visible
	canSeeVotes := isOwner || sess.Visibility != ws.ViewHidden

	options := make([]map[string]any, 0, len(poll.Content.Options))

	for _, opt := range poll.Content.Options {

		entry := map[string]any{
			"option_id": opt.OptionID,
			"text":      opt.Text,
		}

		if canSeeVotes {
			entry["votes"] = opt.Votes
		}

		options = append(options, entry)
	}

	return map[string]any{
		"poll_id":  poll.PollID,
		"question": poll.Content.Question,
		"options":  options,
		"closed":   poll.State.IsClosed,
		"is_owner": isOwner, // useful for frontend owner UI
	}, nil
}
