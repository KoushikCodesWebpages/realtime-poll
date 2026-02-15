package services

import (
	"time"
	"context"
	"errors"
	"github.com/google/uuid"

	"realtime-poll/internal/models"
	"realtime-poll/internal/dto"
	"realtime-poll/internal/repository"
)



type PollCreateService struct{}

func (s *PollCreateService) CreatePoll(
	ctx context.Context,
	userID string,
	userEmail string,
	req dto.CreatePollReq,
) (*models.Poll, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	now := time.Now()

	// validate time window
	if req.StartAt != nil && req.EndAt != nil {
		if req.EndAt.Before(*req.StartAt) {
			return nil, errors.New("end_at cannot be before start_at")
		}
	}

	// build options
	options := make([]models.Option, 0, len(req.Options))
	for _, opt := range req.Options {
		options = append(options, models.Option{
			OptionID: uuid.NewString(),
			Text:     opt,
			Votes:    0,
		})
	}

	// generate share id only for link polls
	shareID := ""
	if req.Visibility == "link" {
		shareID = uuid.NewString()[:8]
	}

	poll := &models.Poll{
		PollID:  uuid.NewString(),
		OwnerID: userID,

		Content: models.ContentSettings{
			Question:    req.Question,
			Description: req.Description,
			Options:     options,
		},

		Access: models.AccessSettings{
			Visibility:    req.Visibility,
			AllowedEmails: req.AllowedEmails,
			RequireLogin:  req.Visibility != "public",
		},

		Vote: models.VoteSettings{
			MaxVotesPerUser: 1,
			AllowChangeVote: req.AllowChange,
			AnonymousVote:   req.Anonymous,
			HideResults:     false,
			ShowVoters:      false,
			UniqueIP:        true,
			UniqueSession:   true,
		},

		Distribution: models.DistributionSettings{
			ShareID: shareID,
		},

		Behavior: models.BehaviorSettings{
			StartAt:   req.StartAt,
			EndAt:     req.EndAt,
			AutoClose: true,
		},

		Analytics: models.AnalyticsSettings{
			TrackViews:    true,
			TrackVoters:   true,
			FraudDetection: true,
		},

		State: models.PollState{
			IsClosed: false,
			IsLocked: false,
			Version:  1,
		},

		Meta: models.Meta{
			CreatedAt:  now,
			UpdatedAt:  now,
			ExpiresAt:  req.EndAt,
			TotalVotes: 0,
			TotalViews: 0,
		},
	}

	if err := repository.CreatePoll(ctx, poll); err != nil {
		return nil, err
	}

	return poll, nil
}
