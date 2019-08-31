package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		randomString, err := makeRandomString(i)
		assert.NoError(t, err)
		assert.Equal(t, i, len(randomString))
	}
}
