package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"
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

func GenerateShareToken(pollID string, mode string, duration time.Duration) (string, error) {

	payload := SharePayload{
		PollID: pollID,
		Mode:   mode,
		Exp:    time.Now().Add(duration).Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(data)

	sig := sign(payloadB64)

	token := payloadB64 + "." + sig
	return token, nil
}

func ParseShareToken(token string) (*SharePayload, error) {

	parts := split(token)
	if len(parts) != 2 {
		return nil, errors.New("malformed token")
	}

	payloadB64 := parts[0]
	sig := parts[1]

	if !verify(payloadB64, sig) {
		return nil, errors.New("invalid signature")
	}

	data, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	var payload SharePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, errors.New("invalid payload")
	}

	if time.Now().Unix() > payload.Exp {
		return nil, errors.New("token expired")
	}

	return &payload, nil
}

func sign(payload string) string {
	h := hmac.New(sha256.New, shareSecret)
	h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func verify(payload, sig string) bool {
	expected := sign(payload)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func split(token string) []string {
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return nil
}

