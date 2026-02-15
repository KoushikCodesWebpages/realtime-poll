package services

import (
	"context"
	"errors"
	"time"

	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
)

type VoteService struct{}

func (s *VoteService) CastVote(
	ctx context.Context,
	pollID string,
	optionID string,
	userID string,
	sessionID string,
	ip string,
) error {

	// -------------------- load poll --------------------
	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil || poll == nil {
		return errors.New("poll not found")
	}

	// -------------------- lifecycle --------------------
	if poll.State.IsClosed {
		return errors.New("poll closed")
	}

	now := time.Now()

	if poll.Behavior.StartAt != nil && now.Before(*poll.Behavior.StartAt) {
		return errors.New("poll not started")
	}

	if poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt) {
		return errors.New("poll ended")
	}

	// -------------------- access rules --------------------

	if poll.Access.RequireLogin && userID == "" {
		return errors.New("login required")
	}

	for _, b := range poll.Access.BlockedUsers {
		if b == userID {
			return errors.New("blocked from voting")
		}
	}

	if poll.Access.Visibility == "whitelist" {
		allowed := false
		for _, u := range poll.Access.AllowedUsers {
			if u == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			return errors.New("not allowed to vote")
		}
	}

	if !poll.Vote.AnonymousVote && userID == "" {
		return errors.New("anonymous voting disabled")
	}

	// -------------------- validate option --------------------
	valid := false
	for _, opt := range poll.Content.Options {
		if opt.OptionID == optionID {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid option")
	}

	// -------------------- vote tracking --------------------
	existing, _ := repository.FindExistingVote(ctx, pollID, userID, sessionID)

	count, _ := repository.CountVotesByUser(ctx, pollID, userID, sessionID, ip)

	// ---------- vote limit ----------
	if poll.Vote.MaxVotesPerUser > 0 && int(count) >= poll.Vote.MaxVotesPerUser && existing == nil {
		return errors.New("vote limit reached")
	}

	// ---------- unique session ----------
	if poll.Vote.UniqueSession && existing == nil && sessionID != "" && count > 0 {
		return errors.New("already voted from this session")
	}

	// ---------- unique ip ----------
	if poll.Vote.UniqueIP && existing == nil && ip != "" && count > 0 {
		return errors.New("already voted from this network")
	}

	// -------------------- first vote --------------------
	if existing == nil {

		vote := models.Vote{
			PollID:    pollID,
			OptionID:  optionID,
			UserID:    userID,
			SessionID: sessionID,
			IP:        ip,
			CreatedAt: time.Now(),
		}

		if err := repository.InsertVote(ctx, vote); err != nil {
			return err
		}

		return repository.IncrementOptionVote(ctx, pollID, optionID, 1)
	}

	// -------------------- change vote --------------------
	if !poll.Vote.AllowChangeVote {
		return errors.New("already voted")
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

	return repository.UpdateVote(ctx, existing.ID, optionID)
}
