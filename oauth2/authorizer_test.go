package oauth2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCodeVerifier(t *testing.T) {
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
