package services

import (
	"context"
	"time"

	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
)

type ViewerState struct {
	CanVote        bool   `json:"can_vote"`
	CanViewResults bool   `json:"can_view_results"`
	VotedOptionID  string `json:"voted_option_id,omitempty"`
	Started        bool   `json:"started"`
	Ended          bool   `json:"ended"`
}

func BuildViewerState(ctx context.Context, poll *models.Poll, userID, sessionID string) (*ViewerState, error) {

	now := time.Now()

	started := poll.Behavior.StartAt == nil || now.After(*poll.Behavior.StartAt)
	ended := poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt)

	canVote := started && !ended

	if poll.Access.Visibility == "authenticated" && userID == "" {
		canVote = false
	}

	if poll.Access.Visibility == "whitelist" && userID == "" {
		canVote = false
	}

	vote, _ := repository.GetVoteForViewer(ctx, poll.PollID, userID, sessionID)

	voted := ""
	if vote != nil {
		voted = vote.OptionID
		if !poll.Vote.AllowChangeVote {
			canVote = false
		}
	}

	canViewResults := !poll.Vote.HideResults || ended || poll.Behavior.ShowLiveResults

	return &ViewerState{
		CanVote:        canVote,
		CanViewResults: canViewResults,
		VotedOptionID:  voted,
		Started:        started,
		Ended:          ended,
	}, nil
}
