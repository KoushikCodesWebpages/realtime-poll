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


type VoteService struct{}

func (s *VoteService) CastVote(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) (int64, error) {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return 0, apperror.Internal()
	}
	if poll == nil {
		return 0, &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	now := time.Now()

	// ===== Time rules =====
	if poll.Behavior.StartAt != nil && now.Before(*poll.Behavior.StartAt) {
		return 0, &apperror.AppError{
			Code:    apperror.POLL_NOT_STARTED,
			Message: "Voting has not started yet",
		}
	}

	if poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt) {
		return 0, &apperror.AppError{
			Code:    apperror.POLL_ENDED,
			Message: "Voting has ended",
		}
	}

	if poll.State.IsClosed {
		return 0, &apperror.AppError{
			Code:    apperror.POLL_CLOSED,
			Message: "Poll is closed",
		}
	}

	// ===== Permission rules =====
	if err := canUserVote(poll, userID); err != nil {
		return 0, err
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
		return 0, &apperror.AppError{
			Code:    apperror.VALIDATION_FAILED,
			Message: "Invalid option",
		}
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

	// ===== Try insert =====
	err = repository.InsertVoteAtomic(ctx, vote)

	// --- First vote ---
	if err == nil {
		version, err := repository.IncrementOptionVote(ctx, pollID, optionID, 1)
		if err != nil {
			return 0, apperror.Internal()
		}
		return version, nil
	}

	// ===== Not duplicate => infra failure =====
	if !repository.IsDuplicateKey(err) {
		return 0, apperror.Internal()
	}

	// ===== Already voted =====
	if !poll.Vote.AllowChangeVote {
		return 0, &apperror.AppError{
			Code:    apperror.POLL_ALREADY_VOTED,
			Message: "You already voted",
		}
	}

	// ===== Change vote =====
	existing, err := repository.GetVoteByIdentity(ctx, pollID, identity)
	if err != nil {
		return 0, apperror.Internal()
	}
	if existing == nil {
		return 0, apperror.Internal()
	}

	// same option = no-op (no new version)
	if existing.OptionID == optionID {
		return poll.State.Version, nil
	}

	// decrement old
	if _, err := repository.IncrementOptionVote(ctx, pollID, existing.OptionID, -1); err != nil {
		return 0, apperror.Internal()
	}

	// increment new → authoritative version
	version, err := repository.IncrementOptionVote(ctx, pollID, optionID, 1)
	if err != nil {
		return 0, apperror.Internal()
	}

	if err := repository.UpdateVoteOption(ctx, existing.VoteID, optionID); err != nil {
		return 0, apperror.Internal()
	}

	return version, nil
}

func (s *VoteService) CastVoteRealtime(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) (int64, error){
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
