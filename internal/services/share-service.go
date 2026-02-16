package services

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"time"
	"strings"

	"crypto/hmac"
	"crypto/sha256"

	"realtime-poll/internal/models"
	"realtime-poll/internal/apperror"
)

func sign(data []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}

func GenerateShareToken(pollID string, mode string, minutes int64) (string, error) {

	claims := models.ShareTokenClaims{
		PollID: pollID,
		Mode:   mode,
	}

	if mode == "timed" {
		claims.Exp = time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
	}

	payload, _ := json.Marshal(claims)

	secret := []byte(os.Getenv("SHARE_TOKEN_SECRET"))
	sig := sign(payload, secret)

	token := base64.RawURLEncoding.EncodeToString(payload) + "." +
		base64.RawURLEncoding.EncodeToString(sig)

	return token, nil
}

func VerifyShareToken(token string) (*models.ShareTokenClaims, error) {

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, apperror.New(apperror.TokenMalformed, "Malformed share token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, apperror.New(apperror.TokenMalformed, "Invalid payload encoding")
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, apperror.New(apperror.TokenMalformed, "Invalid signature encoding")
	}

	secret := []byte(os.Getenv("SHARE_TOKEN_SECRET"))
	expected := sign(payload, secret)

	if !hmac.Equal(sig, expected) {
		return nil, apperror.New(apperror.TokenInvalid, "Signature verification failed")
	}

	var claims models.ShareTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, apperror.New(apperror.TokenMalformed, "Invalid token data")
	}

	if claims.Mode == "timed" && time.Now().Unix() > claims.Exp {
		return nil, apperror.New(apperror.TokenExpired, "Share link expired")
	}

	return &claims, nil
}

