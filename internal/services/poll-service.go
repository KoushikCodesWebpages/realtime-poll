package services

import (
	"time"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
)

func CreatePoll(question string, options []string) (*models.Poll, error) {

	var opts []models.Option
	for _, o := range options {
		opts = append(opts, models.Option{
			ID:    uuid.NewString(),
			Text:  o,
			Votes: 0,
		})
	}

	poll := models.Poll{
		ID:        uuid.NewString(),
		Question:  question,
		Options:   opts,
		CreatedAt: time.Now(),
		IsClosed:  false,
	}

	return &poll, repository.CreatePoll(poll)
}

func Vote(pollID, optionID string) error {
	return repository.IncrementVote(pollID, optionID)
}
