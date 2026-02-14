package services

import (
	"time"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
)

func Register(username, password string) error {

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	return repository.CreateUser(user)
}

func Login(username, password string) (*models.Session, error) {

	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, err
	}

	session := models.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return &session, repository.CreateSession(session)
}
