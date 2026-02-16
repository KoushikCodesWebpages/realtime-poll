package services

import (


	"time"

	"context"


	"realtime-poll/internal/apperror"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
	"realtime-poll/internal/models"
)


type ShareService struct{}

func NewShareService() *ShareService {
	return &ShareService{}
}
func (s *ShareService) GenerateShareLink(
	ctx context.Context,
	userID string,
	pollID string,
) (string, string, error) {

	if userID == "" {
		return "", "", apperror.Unauthorized()
	}

	poll, err := repository.GetPollByID(ctx, pollID)
	if err != nil {
		return "", "", apperror.Internal()
	}
	if poll == nil {
		return "", "", &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	if poll.OwnerID != userID {
		return "", "", &apperror.AppError{
			Code:    apperror.USER_FORBIDDEN,
			Message: "You are not the owner of this poll",
		}
	}

	duration := 365 * 24 * time.Hour

	token, err := utils.GenerateShareToken(
		poll.PollID,
		poll.Access.Visibility,
		int64(duration.Minutes()),
	)
	if err != nil {
		return "", "", apperror.Internal()
	}

	return token, poll.Access.Visibility, nil
}

type SharedPollResult struct {
	Poll   *models.Poll
	UserID string
	SessionID string
	Viewer struct {
		CanVote        bool
		CanViewResults bool
		VotedOptionID  string
		Started        bool
		Ended          bool
	}
}

func (s *ShareService) ViewSharedPoll(
	ctx context.Context,
	token string,
	ip string,
	cookieSessionID string,
) (*SharedPollResult, error) {

	// -------- verify token --------
	claims, err := utils.VerifyShareToken(token)
	if err != nil {
		return nil, err
	}

	// -------- load poll --------
	poll, err := repository.GetPollByID(ctx, claims.PollID)
	if err != nil {
		return nil, apperror.Internal()
	}
	if poll == nil || poll.State.IsDeleted {
		return nil, &apperror.AppError{
			Code:    apperror.POLL_NOT_FOUND,
			Message: "Poll not found",
		}
	}

	// -------- resolve session --------
	var userID string
	sessionID := cookieSessionID

	if sessionID != "" {
		session, _ := repository.GetSession(sessionID)
		if session != nil {
			userID = session.UserID
		}
	}

	// -------- whitelist rule --------
	if claims.Mode == "whitelist" && userID == "" {
		return nil, apperror.Unauthorized()
	}

	// -------- ensure session --------
	if sessionID == "" {
		anon, err := repository.CreateAnonymousSession(ctx, ip)
		if err != nil {
			return nil, apperror.Internal()
		}
		sessionID = anon.SessionID
	}

	// -------- viewer state --------
	viewer, _ := BuildViewerState(ctx, poll, userID, sessionID)

	result := &SharedPollResult{
		Poll:      poll,
		UserID:    userID,
		SessionID: sessionID,
	}

	result.Viewer = struct {
		CanVote        bool
		CanViewResults bool
		VotedOptionID  string
		Started        bool
		Ended          bool
	}{
		viewer.CanVote,
		viewer.CanViewResults,
		viewer.VotedOptionID,
		viewer.Started,
		viewer.Ended,
	}

	return result, nil
}
