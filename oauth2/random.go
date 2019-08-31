package oauth2

import (
	"math/rand"

	"golang.org/x/xerrors"
)

// Generate random values for OAuth2 authorization.

const (
	// Constants used when constructing random strings.

	// Character set that works for URL arguments.
	unreservedCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
)

// makeRandomString returns a random string of the specified length.
func makeRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", xerrors.Errorf("get random byte: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = unreservedCharacters[b%byte(len(unreservedCharacters))]
	}

	return string(bytes), nil
}
