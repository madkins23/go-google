package oauth2

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

// Retrieves a token from a local file and decodes it.
func loadToken(absPath string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, xerrors.Errorf("open token file %s for read: %w", absPath, err)
	}
	defer func() { _ = f.Close() }() // don't care about close when reading

	if err = json.NewDecoder(f).Decode(tok); err != nil {
		err = xerrors.Errorf("reading token from %s: %w", absPath, err)
	}

	return tok, err
}

// Saves a token to a file path.
func saveToken(absPath string, token *oauth2.Token) error {
	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return xerrors.Errorf("open token file %s for write: %w", absPath, err)
	}

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		err = xerrors.Errorf("write token to %s: %w", absPath, err)
	}

	if err := f.Close(); err != nil {
		err = xerrors.Errorf("close token file %s after write: %v", absPath, err)
	}

	return err
}
