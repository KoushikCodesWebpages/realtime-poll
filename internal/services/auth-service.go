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

type LoginResult struct {
	Session *models.Session
	User    *models.User
}

func Login(identifier, password string) (*LoginResult, error) {

	user, err := repository.FindUserByIdentifier(identifier)
	if err != nil || user == nil {
		println("LOGIN FAIL: user not found for", identifier)
		return nil, ErrInvalidCredentials
	}

	println("FOUND USER:", user.Username)
	println("HASH:", user.PasswordHash)

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	session, err := CreateSession(user.AuthUserID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Session: session,
		User:    user,
	}, nil
}


const SessionDuration = 7 * 24 * time.Hour // 7 days

func CreateSession(userID string) (*models.Session, error) {

	now := time.Now().UTC()

	session := &models.Session{
		SessionID: uuid.NewString(),
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionDuration),
	}

	err := repository.InsertSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
