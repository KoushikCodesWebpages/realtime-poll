package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserExists = errors.New("username already exists")

func Register(username, email, password string) error {

	// check existing user
	existing, _ := repository.FindUserByUsername(username)
	if existing != nil {
		return ErrUserExists
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	return repository.CreateUser(user)
}

func Login(identifier, password string) (*models.Session, error) {

	user, err := repository.FindUserByIdentifier(identifier)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	session := models.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := repository.CreateSession(session); err != nil {
		return nil, err
	}

	return &session, nil
}
