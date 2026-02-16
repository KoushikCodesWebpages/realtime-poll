package services

import (
	"context"
	// "errors"
	"time"


	"github.com/google/uuid"

	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/apperror"
	"realtime-poll/internal/ws"
)
func canUserVote(poll *models.Poll, userID string) error {

	switch poll.Access.Visibility {

	case "authenticated":
		if userID == "" {
			return apperror.New(apperror.LoginRequired, "login required to vote")
		}

	case "whitelist":
		if userID == "" {
			return apperror.New(apperror.LoginRequired, "login required to vote")
		}

		allowed := false
		for _, u := range poll.Access.AllowedUsers {
			if u == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			return apperror.New(apperror.NotWhitelisted, "not allowed to vote")
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
) error {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil {
		return apperror.New(apperror.PollNotFound, "poll not found")
	}

	now := time.Now()

	if poll.Behavior.StartAt != nil && now.Before(*poll.Behavior.StartAt) {
		return apperror.New(apperror.PollNotStarted, "poll not started")
	}

	if poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt) {
		return apperror.New(apperror.PollEnded, "poll ended")
	}

	if poll.State.IsClosed {
		return apperror.New(apperror.PollClosed, "poll closed")
	}

	if err := canUserVote(poll, userID); err != nil {
		return err
	}

	valid := false
	for _, opt := range poll.Content.Options {
		if opt.OptionID == optionID {
			valid = true
			break
		}
	}
	if !valid {
		return apperror.New(apperror.VoteInvalidOption, "invalid option")
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

	err = repository.InsertVoteAtomic(ctx, vote)

	if err == nil {
		return repository.IncrementOptionVote(ctx, pollID, optionID, 1)
	}

	if !repository.IsDuplicateKey(err) {
		return err
	}

	if !poll.Vote.AllowChangeVote {
		return apperror.New(apperror.VoteAlreadyCast, "already voted")
	}

	existing, err := repository.GetVoteByIdentity(ctx, pollID, identity)
	if err != nil || existing == nil {
		return apperror.New(apperror.VoteNotAllowed, "vote not found")
	}

	if existing.OptionID == optionID {
		return nil
	}

	if err := repository.IncrementOptionVote(ctx, pollID, existing.OptionID, -1); err != nil {
		return err
	}

	if err := repository.IncrementOptionVote(ctx, pollID, optionID, 1); err != nil {
		return err
	}

	return repository.UpdateVoteOption(ctx, existing.VoteID, optionID)
}

func (s *VoteService) CastVoteRealtime(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) error {
	return s.CastVote(ctx, pollID, optionID, userID, sessionID, ip)
}


func (s *VoteService) GetResults(ctx context.Context, pollID string) ([]ws.OptionResult, error) {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil {
		return nil, apperror.New(apperror.PollNotFound, "poll not found")
	}

	results := make([]ws.OptionResult, 0, len(poll.Content.Options))

	for _, opt := range poll.Content.Options {
		results = append(results, ws.OptionResult{
			OptionID: opt.OptionID,
			Votes:    opt.Votes,
		})
	}

	return results, nil
}
