package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

const (
	driveSecretPath = ".ssh/drive-secret.json"
	driveTokenPath  = ".ssh/drive-token.json"
)

// GetClient acquires a read token and generates an HTTP client for access to the Google API.
func GetClient(scope ...string) (*http.Client, error) {
	// Get the client secret data.
	bytes, err := ioutil.ReadFile(homePath(driveSecretPath))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to read ~/%s: %v", driveSecretPath, err.Error()))
	}

	config, err := google.ConfigFromJSON(bytes, scope...)
	// NOTE: If modifying the scope, delete your previously saved client_secret.json.
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to parse ~/%s: %v", driveSecretPath, err))
	}

	absPath := homePath(driveTokenPath)

	// Attempt to read token.
	tok, err := readToken(absPath)
	if err != nil {
		// Couldn't read token, get one from the web and save it.
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get token from web")
		}

		err = saveToken(absPath, tok)
		if err != nil {
			return nil, errors.Wrap(err, "unable to save token")
		}
	}

	return config.Client(context.Background(), tok), nil
}

// Attempt to verify that Drive is running.
func ping(service *drive.Service) (bool, error) {
	// Test the connection by listing files.
	rows, err := service.Files.List(). /*Fields("id, description").*/ Do()
	if err != nil {
		return false, err
	}

	return len(rows.Items) > 0, nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrap(err, "unable to read authorization code")
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve token from web")
	}

	return tok, nil
}

// Return a path constructed from the specified relative path and the user's home directory.
func homePath(relPath string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to get current user: %v", err)
	}

	return filepath.Join(usr.HomeDir, relPath)
}

// Retrieves a token from a local file and decodes it.
func readToken(absPath string) (*oauth2.Token, error) {
	f, err := os.Open(absPath)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

// Saves a token to a file path.
func saveToken(absPath string, token *oauth2.Token) error {
	fmt.Printf("saving credential file to: %s\n", absPath)

	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to open %s for writing", absPath))
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
