package services

import (
	"time"
	"context"
	// "errors"
	"strings"
	"net/url"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"

	"realtime-poll/internal/utils"
	"realtime-poll/internal/models"
	"realtime-poll/internal/dto"
	"realtime-poll/internal/apperror"
	"realtime-poll/internal/repository"

	
)


type PollEditService struct{}


func (s *PollEditService) PatchPoll(ctx context.Context, userID, pollID string, update bson.M) error {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return apperror.Internal()
	}
	if poll == nil {
		return &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	if poll.OwnerID != userID {
		return &apperror.AppError{
			Code:    apperror.USER_FORBIDDEN,
			Message: "You are not the owner of this poll",
		}
	}

	if poll.State.IsLocked || poll.State.IsClosed {
		return &apperror.AppError{
			Code:    apperror.POLL_EDIT_NOT_ALLOWED,
			Message: "Poll cannot be edited",
		}
	}

	if poll.Meta.TotalVotes > 0 {
		return &apperror.AppError{
			Code:    apperror.POLL_VOTING_STARTED,
			Message: "Poll already has votes",
		}
	}

	update["meta.updated_at"] = time.Now()

	if err := repository.UpdatePollFields(ctx, pollID, update); err != nil {
		return apperror.Internal()
	}

	return nil
}

func (s *PollEditService) PutPoll(ctx context.Context, userID, pollID string, update bson.M) error {

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return apperror.Internal()
	}
	if poll == nil {
		return &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	if poll.OwnerID != userID {
		return &apperror.AppError{
			Code:    apperror.USER_FORBIDDEN,
			Message: "You are not the owner of this poll",
		}
	}

	if poll.Meta.TotalVotes > 0 {
		return &apperror.AppError{
			Code:    apperror.POLL_VOTING_STARTED,
			Message: "Poll already has votes",
		}
	}

	if poll.State.IsLocked || poll.State.IsClosed {
		return &apperror.AppError{
			Code:    apperror.POLL_EDIT_NOT_ALLOWED,
			Message: "Poll cannot be replaced",
		}
	}

	update["meta.updated_at"] = time.Now()

	if err := repository.UpdatePollFields(ctx, pollID, update); err != nil {
		return apperror.Internal()
	}

	return nil
}

func (s *PollEditService) DeletePoll(
	ctx context.Context,
	userID string,
	pollID string,
	force bool,
) error {

	poll, err := repository.GetPollByIDRaw(ctx, pollID)
	if err != nil {
		return apperror.Internal()
	}
	if poll == nil {
		return &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	if poll.OwnerID != userID {
		return &apperror.AppError{
			Code:    apperror.USER_FORBIDDEN,
			Message: "You are not allowed to delete this poll",
		}
	}

	if force {
		if !poll.State.IsDeleted {
			return &apperror.AppError{
				Code:    apperror.VALIDATION_FAILED,
				Message: "Poll must be soft deleted first",
			}
		}
		if err := repository.HardDeletePoll(ctx, pollID); err != nil {
			return apperror.Internal()
		}
		return nil
	}

	if poll.State.IsDeleted {
		return &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll already deleted",
		}
	}

	if err := repository.SoftDeletePoll(ctx, pollID); err != nil {
		return apperror.Internal()
	}

	return nil
}




type PollGetService struct{}
func (s *PollGetService) GetMyPollsPaginated(
	ctx context.Context,
	userID string,
	filter dto.PollFilter,
	query url.Values,
	cursor *time.Time,
	direction string,
) (*dto.PollListResponse, error) {

	if userID == "" {
		return nil, apperror.Unauthorized()
	}

	polls, err := repository.GetPollsByOwnerPaginated(ctx, userID, filter, cursor, direction)
	if err != nil {
		return nil, apperror.Internal()
	}

	total, err := repository.CountPollsByOwner(ctx, userID)
	if err != nil {
		return nil, apperror.Internal()
	}

	base := "/b1/poll/mine"

	var next *string
	var prev *string

	if len(polls) > 0 {

		if int64(len(polls)) == filter.Limit {
			n := utils.BuildPageLink(base, query, polls[len(polls)-1].Meta.CreatedAt, "next")
			next = &n
		}

		if cursor != nil {
			p := utils.BuildPageLink(base, query, polls[0].Meta.CreatedAt, "prev")
			prev = &p
		}
	}

	return &dto.PollListResponse{
		Data:  polls,
		Next:  next,
		Prev:  prev,
		Total: total,
	}, nil
}

func (s *PollGetService) GetMyPolls(ctx context.Context, userID string) ([]models.Poll, error) {

	if userID == "" {
		return nil, apperror.Unauthorized()
	}

	polls, err := repository.GetPollsByOwner(ctx, userID)
	if err != nil {
		return nil, apperror.Internal()
	}

	return polls, nil
}


type PollSingleService struct{}

func (s *PollSingleService) GetPoll(
	ctx context.Context,
	userID string,
	pollID string,
) (*models.Poll, error) {

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

	// owner only (for now)
	if poll.OwnerID != userID {
		return nil, &apperror.AppError{
			Code:    apperror.USER_FORBIDDEN,
			Message: "You are not allowed to access this poll",
		}
	}

	// auto close expired
	if poll.Meta.ExpiresAt != nil && time.Now().After(*poll.Meta.ExpiresAt) {
		poll.State.IsClosed = true
	}

	return poll, nil
}
type PollCreateService struct{}

func (s *PollCreateService) CreatePoll(
	ctx context.Context,
	userID string,
	userEmail string,
	req dto.CreatePollReq,
) (*models.Poll, error) {

	// ---------- AUTH ----------
	if userID == "" {
		return nil, apperror.Unauthorized()
	}

	now := time.Now().UTC()

	// ---------- TIME VALIDATION ----------
	if req.StartAt != nil && req.EndAt != nil {
		if req.EndAt.Before(*req.StartAt) {
			return nil, &apperror.AppError{
				Code:    apperror.VALIDATION_FAILED,
				Message: "end_at cannot be before start_at",
			}
		}
	}

	// ---------- OPTIONS ----------
	if len(req.Options) < 2 {
		return nil, apperror.Validation("minimum 2 options required")
	}

	options := make([]models.Option, 0, len(req.Options))
	for _, opt := range req.Options {
		options = append(options, models.Option{
			OptionID: uuid.NewString(),
			Text:     strings.TrimSpace(opt),
			Votes:    0,
		})
	}

	// ---------- VISIBILITY ----------
	requireLogin := req.Visibility != "public"

	// link polls get share id
	shareID := ""
	if req.Visibility == "link" {
		shareID = uuid.NewString()[:8]
	}

	// ---------- BUILD POLL ----------
	poll := &models.Poll{
		PollID:  uuid.NewString(),
		OwnerID: userID,

		// ================= CONTENT =================
		Content: models.ContentSettings{
			Question:    req.Question,
			Description: req.Description,
			Options:     options,

			Images:      req.Images,
			AllowCustom: req.AllowCustomOption,
			Randomize:   req.RandomizeOptions,
		},

		// ================= ACCESS =================
		Access: models.AccessSettings{
			Visibility:    req.Visibility,
			AllowedEmails: req.AllowedEmails,
			RequireLogin:  requireLogin,
		},

		// ================= VOTE =================
		Vote: models.VoteSettings{
			MaxVotesPerUser: req.MaxVotesPerUser,
			AllowChangeVote: req.AllowChange,
			AnonymousVote:   req.Anonymous,
			HideResults:     req.HideResultsUntilEnd,
			ShowVoters:      req.ShowVoters,

			UniqueIP:      req.UniqueIP,
			UniqueSession: req.UniqueSession,
		},

		// ================= DISTRIBUTION =================
		Distribution: models.DistributionSettings{
			ShareID: shareID,
		},

		// ================= BEHAVIOR =================
		Behavior: models.BehaviorSettings{
			StartAt:          req.StartAt,
			EndAt:            req.EndAt,
			AutoClose:        req.AutoClose,
			ShowLiveResults:  req.ShowLiveResults,
			NotifyOwner:      req.NotifyOwnerOnVote,
		},

		// ================= ANALYTICS =================
		Analytics: models.AnalyticsSettings{
			TrackViews:     true,
			TrackVoters:    true,
			TrackLocation:  false,
			TrackDevice:    false,
			FraudDetection: true,
		},

		// ================= STATE =================
		State: models.PollState{
			IsClosed:  false,
			IsLocked:  false,
			IsDeleted: false,
			Version:   1,
		},

		// ================= META =================
		Meta: models.Meta{
			CreatedAt:  now,
			UpdatedAt:  now,
			LastVote:   nil,
			TotalVotes: 0,
			TotalViews: 0,
			ExpiresAt:  req.EndAt, // derived from behavior
		},
	}

	// ---------- INSERT ----------
	if err := repository.CreatePoll(ctx, poll); err != nil {
		return nil, apperror.Internal()
	}

	return poll, nil
}
