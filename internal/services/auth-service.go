package services

import (

	"time"
	"strings"

	"github.com/google/uuid"
	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
	"realtime-poll/internal/utils"
	"realtime-poll/internal/apperror"
)




func Register(username, email, password string) error {

	username = strings.ToLower(strings.TrimSpace(username))
	email = strings.ToLower(strings.TrimSpace(email))

	// ===== Username check =====
	existingUser, err := repository.FindUserByUsername(username)
	if err != nil {
		return apperror.Internal()
	}
	if existingUser != nil {
		return &apperror.AppError{
			Code:    apperror.USER_EXISTS,
			Message: "Username already exists",
		}
	}

	// ===== Email check =====
	existingEmail, err := repository.FindUserByEmail(email)
	if err != nil {
		return apperror.Internal()
	}
	if existingEmail != nil {
		return &apperror.AppError{
			Code:    apperror.USER_EXISTS,
			Message: "Email already registered",
		}
	}

	// ===== Password hash =====
	hash, err := utils.HashPassword(password)
	if err != nil {
		return apperror.Internal()
	}

	user := models.User{
		AuthUserID:   uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	// ===== Insert =====
	if err := repository.CreateUser(user); err != nil {
		return apperror.Internal()
	}

	return nil
}


type LoginResult struct {
	Session *models.Session
	User    *models.User
}
func Login(identifier, password string) (*LoginResult, error) {

	identifier = strings.ToLower(strings.TrimSpace(identifier))

	// ===== Find user =====
	user, err := repository.FindUserByIdentifier(identifier)
	if err != nil {
		return nil, apperror.Internal()
	}
	if user == nil {
		return nil, apperror.InvalidCredentials()
	}

	// ===== Password check =====
	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, apperror.InvalidCredentials()
	}

	// ===== Session create =====
	session, err := CreateSession(user.AuthUserID)
	if err != nil {
		return nil, apperror.Internal()
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

	if err := repository.InsertSession(session); err != nil {
		return nil, apperror.Internal()
	}

	return session, nil
}
