package slack

import (
	"crypto/rand"
	"encoding/base64"
)

const OAuthStateBytes = 32

func GenerateState() string {
	b := make([]byte, OAuthStateBytes)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure random state: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}
