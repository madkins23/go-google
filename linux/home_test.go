package linux

import (
	"os/user"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePath(t *testing.T) {
	relPath := "relative/path.ext"
	homePath, err := HomePath(relPath)
	assert.NoError(t, err)
	current, err := user.Current()
	assert.NoError(t, err)
	assert.Equal(t, path.Join(current.HomeDir, relPath), homePath)
}
