package services

import (
	"context"

	"realtime-poll/internal/apperror"
	"realtime-poll/internal/dto"
	"realtime-poll/internal/repository"
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

