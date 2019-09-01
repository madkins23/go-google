package oauth2

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/browser"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"

	"github.com/madkins23/go-utils/path"
)

// Authorizer provides a client for accessing Google APIs.
// Interacts with Google OAuth2 endpoints to acquire authorization as necessary.
// A secret file must be pre-configured as ~/.ssh/<appName>-secret.json.
type Authorizer interface {
	GetClient() (*http.Client, error)
}

// Authorization data including loopback server listener.
type authorizer struct {
	appName   string
	port      string
	state     string
	scopes    []string
	config    *goauth2.Config
	listener  net.Listener
	server    http.Server
	token     *goauth2.Token
	tokenPath string
}

// NewAuthorizer returns a new Authorizer object for the specified application name and access scopes.
// The application name is used in the pathnames for the secret and token files.
// The access scopes specify the required access permissions for the application.
// Access scopes are specified as the final part of the access scope URLs
// (e.g. https://www.googleapis.com/auth/drive is specified as "drive").
// All access scopes must be specified here (no incremental access requests).
func NewAuthorizer(applicationName string, accessScopes []string) (Authorizer, error) {
	var err error

	state, err := makeStateString()
	if err != nil {
		return nil, xerrors.Errorf("make auth state: %w", err)
	}

	authorizer := &authorizer{
		appName: applicationName,
		scopes:  fixedAccessScopes(accessScopes),
		state:   state,
	}

	return authorizer, nil
}

const (
	driveSecretPathFmt = ".ssh/%s-secret.json"
	driveTokenPathFmt  = ".ssh/%s-token.json"
)

// GetClient returns an HTTP client that can be used to access Google APIs.
// Handles authorization for the application.
func (auth *authorizer) GetClient() (*http.Client, error) {
	// Get the client secret data.
	secretPath, err := path.HomePath(fmt.Sprintf(driveSecretPathFmt, auth.appName))
	if err != nil {
		return nil, xerrors.Errorf("get secret path: %w", err)
	}

	bytes, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, xerrors.Errorf("read from secret file: %w", err)
	}

	auth.config, err = google.ConfigFromJSON(bytes, auth.scopes...)
	if err != nil {
		return nil, xerrors.Errorf("config from secret: %w", err)
	}

	auth.tokenPath, err = path.HomePath(fmt.Sprintf(driveTokenPathFmt, auth.appName))
	if err != nil {
		return nil, xerrors.Errorf("get token path: %w", err)
	}

	auth.token, err = loadToken(auth.tokenPath)
	if err != nil {
		fmt.Println("Acquiring token from Google OAuth2 server.")

		// Configure listener to be used for loopback server on a random port.
		auth.listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, xerrors.Errorf("start listener on random port: %w", err)
		}

		// Get port number from listener.
		requestUrl, err := url.ParseRequestURI("http://" + auth.listener.Addr().String())
		if err != nil {
			return nil, xerrors.Errorf("parse server url %s: %w", requestUrl, err)
		}

		if requestUrl.Port() == "" {
			return nil, xerrors.New("empty loopback server port")
		}

		if strings.Contains(requestUrl.Port(), ":") {
			return nil, xerrors.New("loopback server port contains a colon")
		}

		// Configure redirect URL using listener port.
		auth.config.RedirectURL = "http://127.0.0.1:" + requestUrl.Port()

		// Get Google authorization URL and attempt to open a browser tab.
		authUrl := auth.config.AuthCodeURL(auth.state, goauth2.AccessTypeOffline)
		if browser.OpenURL(authUrl) != nil {
			fmt.Printf("Open URL %s in browser", authUrl)
		}

		// Start the server and wait for a callbackResponesFmt.
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			auth.handleCallback(writer, request)
		})

		err = auth.server.Serve(auth.listener)
		if err != nil && strings.HasSuffix(err.Error(), "Server closed") {
			err = nil
		}
	}

	if auth.token == nil {
		return nil, xerrors.New("unable to acquire token")
	}

	return auth.config.Client(context.Background(), auth.token), nil
}

const (
	callbackResponesFmt = `
<html><head><title>Google Authorization for %s</title></head>
<body>
 <h1>Google Authorization for application %s</h1>
 <p>%s</p>
</body>
</html>`
)

// handleCallback
func (auth *authorizer) handleCallback(writer http.ResponseWriter, request *http.Request) {
	message := "Unknown error prior to setting message"

	urlData, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		message = fmt.Sprintf("Error parsing request URI: %v", err)
	}
	query := urlData.Query()

	if state := query.Get("state"); auth.state != state {
		message = fmt.Sprintf("Authorization state does not match: '%s' != '%s'", auth.state, state)
	} else if code := query.Get("code"); code == "" {
		message = "Authorization code is empty"
	} else if auth.token, err = auth.config.Exchange(context.Background(), query.Get("code")); err != nil {
		message = fmt.Sprintf("Unable to get token: %v", err.Error())
	} else if err = saveToken(auth.tokenPath, auth.token); err != nil {
		message = fmt.Sprintf("Unable to save token: %v", err.Error())
	} else {
		message = "Authorization successful"
	}

	_, err = fmt.Fprint(writer, fmt.Sprintf(callbackResponesFmt, auth.appName, auth.appName, message))
	if err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	// Shutdown server after done handling request.
	go func() {
		if err = auth.server.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down server: %v\n", err)
		} else {
			fmt.Printf("Embedded HTTP server closed")
		}
	}()
}

//////////////////////////////////////////////////////////////////////////////
// Utilities

const (
	// Fixed prefix on every access scope string to simplify usage.
	accessScopePrefix = "https://www.googleapis.com/auth/"
)

func fixedAccessScopes(scopes []string) []string {
	fixed := make([]string, len(scopes))

	for index, scope := range scopes {
		fixed[index] = accessScopePrefix + scope
	}

	return fixed
}

const (
	// Authorization request/callbackResponesFmt state size range.
	stateLengthLow  = 15
	stateLengthHigh = 31
)

// Make a random state string.
func makeStateString() (string, error) {
	state, err := makeRandomString(stateLengthLow + rand.Intn(stateLengthHigh-stateLengthLow))

	if err != nil {
		return "", xerrors.Errorf("get random string: %w", err)
	}

	return state, nil
}

const (
	// Authorization challenge size range.
	verifierLengthLow  = 43
	verifierLengthHigh = 128
)

// makeCodeVerifier returns a random code verifier string.
func makeCodeVerifier() (string, error) {
	return makeRandomString(verifierLengthLow + rand.Intn(verifierLengthHigh-verifierLengthLow))
}
