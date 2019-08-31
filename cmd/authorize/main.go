// The authorize application is used to acquire an initial access token for Google API usage.
// Originally this was to be separated from the actual applications to be written,
// but the initial access token acquisition is now in madkins23/oauth2.authorizer.go,
// so now this application is really kind of redundant except for testing.
package main

import (
	"log"
	"os"

	"github.com/madkins23/go-google/oauth2"

	"golang.org/x/xerrors"
)

type clientData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthUri      string `json:"auth_uri"`
	TokenUri     string `json:"token_uri"`
}

type secretData struct {
	Installed clientData `json:"installed"`
}

//////////////////////////////////////////////////////////////////////////////
// Main routine.

func main() {
	// Parse arguments.
	if len(os.Args) < 3 {
		log.Panic("usage:  authorize <applicationName> <accessScope>+")
	}

	// Create Google authorizer.
	authorizer, err := oauth2.NewAuthorizer(os.Args[1], os.Args[2:])
	if err != nil {
		log.Panicf("Unable to create authorizer: %v", err)
	}

	client, err := authorizer.GetClient()
	if err == nil && client == nil {
		err = xerrors.New("no client returned")
	}
	if err != nil {
		log.Panicf("Unable to get client: %v", err)
	}
}
