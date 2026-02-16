package services

import (
	"context"
	// "errors"
	"time"


	"github.com/google/uuid"

	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/apperror"
	"realtime-poll/internal/dto"
)

func canUserVote(poll *models.Poll, userID string) error {

	switch poll.Access.Visibility {

	case "authenticated":
		if userID == "" {
			return apperror.Unauthorized()
		}

	case "whitelist":
		if userID == "" {
			return apperror.Unauthorized()
		}

		allowed := false
		for _, u := range poll.Access.AllowedUsers {
			if u == userID {
				allowed = true
				break
			}
		}

		if !allowed {
			return &apperror.AppError{
				Code:    apperror.USER_FORBIDDEN,
				Message: "You are not allowed to vote in this poll",
			}
		}
	}

	return nil
}

type VoteResult struct {
	Version        int64  `json:"version"`
	SelectedOption string `json:"selected_option"`
	AlreadyVoted   bool   `json:"already_voted"`
	Changed        bool   `json:"changed"`
}

type VoteService struct{}
func (s *VoteService) CastVote(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) (*VoteResult, error) {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, apperror.Internal()
	}
	if poll == nil {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	now := time.Now()

	// ===== Time rules =====
	if poll.Behavior.StartAt != nil && now.Before(*poll.Behavior.StartAt) {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_NOT_STARTED,
			Message: "Voting has not started yet",
		}
	}

	if poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt) {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_ENDED,
			Message: "Voting has ended",
		}
	}

	if poll.State.IsClosed {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_CLOSED,
			Message: "Poll is closed",
		}
	}

	// ===== Permission rules =====
	if err := canUserVote(poll, userID); err != nil {
		return nil, err
	}

	// ===== Validate option =====
	valid := false
	for _, opt := range poll.Content.Options {
		if opt.OptionID == optionID {
			valid = true
			break
		}
	}
	if !valid {
		return nil, apperror.Validation("Invalid option")
	}

	identity := userID
	if identity == "" {
		identity = sessionID
	}

	vote := models.VoteRecord{
		VoteID:    uuid.NewString(),
		PollID:    pollID,
		Identity:  identity,
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ip,
		OptionID:  optionID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// ===== Try insert (first vote) =====
	err = repository.InsertVoteAtomic(ctx, vote)

	if err == nil {
		version, err := repository.IncrementOptionVote(ctx, pollID, optionID, 1)
		if err != nil {
			return nil, apperror.Internal()
		}

		return &VoteResult{
			Version:        version,
			SelectedOption: optionID,
			AlreadyVoted:   false,
			Changed:        false,
		}, nil
	}

	// ===== Infra failure =====
	if !repository.IsDuplicateKey(err) {
		return nil, apperror.Internal()
	}

	// ===== Already voted =====
	existing, err := repository.GetVoteByIdentity(ctx, pollID, identity)
	if err != nil {
		return nil, apperror.Internal()
	}
	if existing == nil {
		return nil, apperror.Internal()
	}

	// ===== No change allowed =====
	if !poll.Vote.AllowChangeVote {
		return &VoteResult{
			Version:        poll.State.Version,
			SelectedOption: existing.OptionID,
			AlreadyVoted:   true,
			Changed:        false,
		}, nil
	}

	// ===== Same option clicked again =====
	if existing.OptionID == optionID {
		return &VoteResult{
			Version:        poll.State.Version,
			SelectedOption: existing.OptionID,
			AlreadyVoted:   true,
			Changed:        false,
		}, nil
	}

	// ===== Change vote =====

	// decrement old
	if _, err := repository.IncrementOptionVote(ctx, pollID, existing.OptionID, -1); err != nil {
		return nil, apperror.Internal()
	}

	// increment new
	version, err := repository.IncrementOptionVote(ctx, pollID, optionID, 1)
	if err != nil {
		return nil, apperror.Internal()
	}

	// update record
	if err := repository.UpdateVoteOption(ctx, existing.VoteID, optionID); err != nil {
		return nil, apperror.Internal()
	}

	return &VoteResult{
		Version:        version,
		SelectedOption: optionID,
		AlreadyVoted:   true,
		Changed:        true,
	}, nil
}


func (s *VoteService) CastVoteRealtime(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) (*VoteResult, error){
	return s.CastVote(ctx, pollID, optionID, userID, sessionID, ip)
}

func (s *VoteService) GetResults(ctx context.Context, pollID string) ([]dto.OptionResult, error) {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, apperror.Internal()
	}

	if poll == nil {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	results := make([]dto.OptionResult, 0, len(poll.Content.Options))

	for _, opt := range poll.Content.Options {
		results = append(results, dto.OptionResult{
			OptionID: opt.OptionID,
			Votes:    opt.Votes,
		})
	}

	return results, nil
}
