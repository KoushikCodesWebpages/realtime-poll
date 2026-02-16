package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	// "errors
	"os"
	"time"
	"strings"


	"realtime-poll/internal/models"
	"realtime-poll/internal/apperror"
)


type SharePayload struct {
	PollID string `json:"pid"`
	Mode   string `json:"mode"` // public | authenticated | whitelist | link
	Exp    int64  `json:"exp"`  // unix timestamp
}

var shareSecret = []byte(getSecret())

func getSecret() string {
	s := os.Getenv("SHARE_TOKEN_SECRET")
	if s == "" {
		return "dev-secret-change-in-prod"
	}
	return s
}

func GenerateShareToken(pollID string, mode string, minutes int64) (string, error) {

	claims := models.ShareTokenClaims{
		PollID: pollID,
		Mode:   mode,
	}

	if mode == "timed" {
		claims.Exp = time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", apperror.Internal()
	}

	secretStr := os.Getenv("SHARE_TOKEN_SECRET")
	if secretStr == "" {
		return "", apperror.Internal()
	}
	secret := []byte(secretStr)

	// sign RAW payload
	signature := sign(payload, secret)

	// encode both parts
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	sigB64 := base64.RawURLEncoding.EncodeToString(signature)

	return payloadB64 + "." + sigB64, nil
}


func sign(data []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}
func verify(payloadB64, sigB64 string, secret []byte) bool {

	payload, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return false
	}

	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return false
	}

	expected := sign(payload, secret)

	return hmac.Equal(sig, expected)
}


func split(token string) []string {
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return nil
}

func VerifyShareToken(token string) (*models.ShareTokenClaims, error) {

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, apperror.TokenInvalid()
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, apperror.TokenInvalid()
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, apperror.TokenInvalid()
	}

	secretStr := os.Getenv("SHARE_TOKEN_SECRET")
	if secretStr == "" {
		return nil, apperror.Internal()
	}
	secret := []byte(secretStr)

	expected := sign(payload, secret)

	if !hmac.Equal(sig, expected) {
		return nil, apperror.TokenInvalid()
	}

	var claims models.ShareTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, apperror.TokenInvalid()
	}

	if claims.Mode == "timed" && time.Now().Unix() > claims.Exp {
		return nil, apperror.TokenExpired()
	}

	return &claims, nil
}
