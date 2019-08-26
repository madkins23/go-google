package linux

import (
	"github.com/pkg/errors"
	"os/user"
	"path/filepath"
)

// Return a path constructed from the specified relative path and the user's home directory.
func HomePath(relPath string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Error getting current user data")
	} else if usr == nil {
		return "", errors.New("Unable to get current user data")
	}

	return filepath.Join(usr.HomeDir, relPath), nil
}
