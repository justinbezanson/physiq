package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const (
	CookieName = "session_token"
	TTL        = 30 * 24 * time.Hour
)

func Generate() (cookieValue string, tokenHash []byte, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, err
	}
	return base64.RawURLEncoding.EncodeToString(raw), hashToken(raw), nil
}

func HashFromCookieValue(cookieValue string) ([]byte, bool) {
	raw, err := base64.RawURLEncoding.DecodeString(cookieValue)
	if err != nil || len(raw) != 32 {
		return nil, false
	}
	return hashToken(raw), true
}

func hashToken(raw []byte) []byte {
	sum := sha256.Sum256(raw)
	return sum[:]
}
