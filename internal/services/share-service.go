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
		return nil, errors.New("invalid token")
	}

	payload, _ := base64.RawURLEncoding.DecodeString(parts[0])
	sig, _ := base64.RawURLEncoding.DecodeString(parts[1])

	secret := []byte(os.Getenv("SHARE_TOKEN_SECRET"))
	expected := sign(payload, secret)

	if !hmac.Equal(sig, expected) {
		return nil, errors.New("signature mismatch")
	}

	var claims models.ShareTokenClaims
	json.Unmarshal(payload, &claims)

	if claims.Mode == "timed" && time.Now().Unix() > claims.Exp {
		return nil, errors.New("link expired")
	}

	return &claims, nil
}
