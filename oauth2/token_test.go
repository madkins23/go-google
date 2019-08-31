package oauth2

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/require"
)

const (
	tokenFile = "token.json"
)

func TestSaveLoad(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "token_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	absPath := path.Join(tempDir, tokenFile)
	token := &oauth2.Token{
		AccessToken:  "access-token",
		TokenType:    "token-type",
		RefreshToken: "refresh-token",
	}

	err = saveToken(absPath, token)
	assert.NoError(t, err)
	loaded, err := loadToken(absPath)
	assert.NoError(t, err)
	assert.Equal(t, token, loaded)
}
