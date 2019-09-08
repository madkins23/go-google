package oauth2

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
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
	ctxt      context.Context
	listener  net.Listener
	server    http.Server
	token     *goauth2.Token
	tokenPath string
	verifier  goauth2.AuthCodeOption
}

// NewAuthorizer returns a new Authorizer object for the specified application name and access scopes.
// The application name is used in the pathnames for the secret and token files.
// The access scopes specify the required access permissions for the application.
// Access scopes are specified as the final part of the access scope URLs
// (e.g. https://www.googleapis.com/auth/drive is specified as "drive").
// All access scopes must be specified here (no incremental access requests).
func NewAuthorizer(applicationName string, accessScopes []string) (Authorizer, error) {
	var err error

	state, err := makeState()
	if err != nil {
		return nil, xerrors.Errorf("make auth state: %w", err)
	}

	authorizer := &authorizer{
		appName: applicationName,
		ctxt:    context.Background(),
		scopes:  fixedAccessScopes(accessScopes),
		state:   state,
	}

	return authorizer, nil
}

const (
	driveSecretPathFmt = ".ssh/%s-secret.json"
	driveTokenPathFmt  = ".ssh/%s-token.json"
)

var (
	challengeMethod = goauth2.SetAuthURLParam("code_challenge_method", "plain")
)

var (
	errNoAuthToken     = xerrors.New("no auth token")
	errEmptyServerPort = xerrors.New("empty loopback server port")
	errServerPortColon = xerrors.New("loopback server port contains a colon")
)

// GetClient returns an HTTP client that can be used to access Google APIs.
// Handles authorization for the application, creating token file if necessary.
// Get service from appropriate API (e.g. google.golang.org/api/drive.NewService()).
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
	if err == nil && auth.token == nil {
		err = errNoAuthToken
	}

	if err != nil {
		fmt.Printf("Acquiring token from Google OAuth2 server because:\n  %v\n", err.Error())

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
			return nil, errEmptyServerPort
		}

		if strings.Contains(requestUrl.Port(), ":") {
			return nil, errServerPortColon
		}

		// Configure redirect URL using listener port.
		auth.config.RedirectURL = "http://127.0.0.1:" + requestUrl.Port()

		// Code verifier to URL.
		codeVerify, err := makeCodeVerifier()
		if err != nil {
			return nil, xerrors.Errorf("get code verifier: %w", err)
		}

		// The code verifier is sent using different keys at different times (WTF?).
		//  Used during authorization call just below.
		verifier := goauth2.SetAuthURLParam("code_challenge", codeVerify)
		//  Used in the loopback server callback when calling golang.org/x/oauth2.Exchange().
		auth.verifier = goauth2.SetAuthURLParam("code_verifier", codeVerify)

		// Get Google authorization URL.
		authUrl := auth.config.AuthCodeURL(auth.state, goauth2.AccessTypeOffline, challengeMethod, verifier)

		fmt.Printf("--> Authorization URL:\n  %s\n", authUrl)

		// Open authentication URL in browser.
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
		return nil, errNoAuthToken
	}

	return auth.config.Client(auth.ctxt, auth.token), nil
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
	urlData, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		err = xerrors.Errorf("parse request URI: %w", err)
	} else if urlData.Path != "/" {
		http.NotFound(writer, request)
		return
	} else if state := urlData.Query().Get("state"); auth.state != state {
		err = xerrors.Errorf("check authorization state (%s != %s): %w", auth.state, state, err)
	} else if code := urlData.Query().Get("code"); code == "" {
		err = xerrors.Errorf("empty authorization code: %w", err)
	} else if auth.token, err = auth.config.Exchange(
		auth.ctxt, code, goauth2.AccessTypeOffline, challengeMethod, auth.verifier); err != nil {
		err = xerrors.Errorf("get token: %w", err)
	} else if auth.token == nil {
		err = errNoAuthToken
	} else if err = saveToken(auth.tokenPath, auth.token); err != nil {
		err = xerrors.Errorf("save token: %w", err)
	}

	message := "Authorization successful, browser tab can be closed."
	if err != nil {
		message = fmt.Sprintf("*** Error! %v\n", err)
		fmt.Fprintln(os.Stderr, message)
	}

	_, err = fmt.Fprint(writer, fmt.Sprintf(callbackResponesFmt, auth.appName, auth.appName, message))
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error writing response: %v\n", err)
	}

	// Shutdown server after done handling request.
	go func() {
		if err = auth.server.Shutdown(auth.ctxt); err != nil {
			fmt.Fprintf(os.Stderr, "*** Error shutting down server: %v\n", err)
		} else {
			fmt.Println("Embedded HTTP server closed")
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

// makeState makes a random state string.
func makeState() (string, error) {
	state, err := makeRandomString(stateLengthLow + rand.Intn(stateLengthHigh-stateLengthLow))

	if err != nil {
		return "", xerrors.Errorf("get random string: %w", err)
	}

	return state, nil
}

const (
	// Authorization verifier size range.
	verifierLengthLow  = 43
	verifierLengthHigh = 128
)

// makeCodeVerifier returns a random code verifier string.
func makeCodeVerifier() (string, error) {
	return makeRandomString(verifierLengthLow + rand.Intn(verifierLengthHigh-verifierLengthLow))
}
