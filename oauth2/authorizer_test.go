package oauth2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizer_fixAccessScopes(t *testing.T) {
	scopes := []string{
		"alpha",
		"bravo",
		"charlie",
	}

	fixed := fixedAccessScopes(scopes)
	assert.Equal(t, "alpha", scopes[0])
	assert.Equal(t, "bravo", scopes[1])
	assert.Equal(t, "charlie", scopes[2])
	assert.Equal(t, "https://www.googleapis.com/auth/alpha", fixed[0])
	assert.Equal(t, "https://www.googleapis.com/auth/bravo", fixed[1])
	assert.Equal(t, "https://www.googleapis.com/auth/charlie", fixed[2])
}

func TestAuthorizer_makeCodeVerifier(t *testing.T) {
	for i := 0; i < 10; i++ {
		codeVerifier, err := makeCodeVerifier()
		assert.NoError(t, err)
		assert.LessOrEqual(t, verifierLengthLow, len(codeVerifier))
		assert.GreaterOrEqual(t, verifierLengthHigh, len(codeVerifier))
		for _, c := range codeVerifier {
			strings.Contains(unreservedCharacters, string(c))
		}
	}
}

func TestAuthorizer_makeState(t *testing.T) {
	for i := 0; i < 10; i++ {
		codeVerifier, err := makeState()
		assert.NoError(t, err)
		assert.LessOrEqual(t, stateLengthLow, len(codeVerifier))
		assert.GreaterOrEqual(t, stateLengthHigh, len(codeVerifier))
		for _, c := range codeVerifier {
			strings.Contains(unreservedCharacters, string(c))
		}
	}
}
