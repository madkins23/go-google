package drive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	client, err := GetClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
