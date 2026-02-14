package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
)

var ErrEmailExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserExists = errors.New("username already exists")
func Register(username, email, password string) error {

	// username check
	existingUser, err := repository.FindUserByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserExists
	}

	// email check
	existingEmail, err := repository.FindUserByEmail(email)
	if err != nil {
		return err
	}
	if existingEmail != nil {
		return ErrEmailExists
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		AuthUserID:   uuid.NewString(),
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
		SessionID:        uuid.NewString(),
		UserID:    user.AuthUserID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := repository.CreateSession(session); err != nil {
		return nil, err
	}

	return &session, nil
}
