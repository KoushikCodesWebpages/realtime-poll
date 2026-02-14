package services

import (
	"time"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
)
func CreatePoll(userID, question string, options []string) (*models.Poll, error) {

	var opts []models.Option
	for _, o := range options {
		opts = append(opts, models.Option{
			OptionID: uuid.NewString(),
			Text:     o,
			Votes:    0,
		})
	}

	poll := models.Poll{
		AuthUserID: userID,          // ← IMPORTANT
		PollID:     uuid.NewString(),
		Question:   question,
		Options:    opts,
		CreatedAt:  time.Now(),
		IsClosed:   false,
	}

	return &poll, repository.CreatePoll(poll)
}


func Vote(pollID, optionID string) error {
	return repository.IncrementVote(pollID, optionID)
}
