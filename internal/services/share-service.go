package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"
	"strings"

	"crypto/hmac"
	"crypto/sha256"

	"realtime-poll/internal/models"
)

func sign(data []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}

func GenerateShareToken(pollID string) (string, error) {

	expMin := os.Getenv("SHARE_TOKEN_EXPIRY_MINUTES")
	expDur, _ := time.ParseDuration(expMin + "m")

	claims := models.ShareTokenClaims{
		PollID: pollID,
		Type:   "view",
		Uses:   50,
		Exp:    time.Now().Add(expDur).Unix(),
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
		return nil, errors.New("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid payload")
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid signature")
	}

	secret := []byte(os.Getenv("SHARE_TOKEN_SECRET"))
	expected := sign(payload, secret)

	if !hmac.Equal(sig, expected) {
		return nil, errors.New("signature mismatch")
	}

	var claims models.ShareTokenClaims
	json.Unmarshal(payload, &claims)

	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}
