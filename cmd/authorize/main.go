// The authorize command acquires an initial access token for Google API usage.
// See the README.md file for more information.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"

	"github.com/madkins23/go-google/oauth2"

	"golang.org/x/xerrors"
)

const (
	usage = `
usage:  authorize <applicationName> <accessScope>*
    <applicatonName> name of application to be authorized
    <accessScope>    access scopes for application (defaults to "drive")
`
)

//////////////////////////////////////////////////////////////////////////////
// Main routine.

func main() {
	fmt.Println("authorize starting")

	// Parse arguments.
	if len(os.Args) < 2 {
		log.Panic(usage)
	}

	applicationName := os.Args[1]

	var accessScopes []string
	if len(os.Args) < 3 {
		accessScopes = []string{"drive"}
	} else {
		accessScopes = os.Args[2:]
	}

	fmt.Printf("Application:   %s\n", applicationName)
	fmt.Printf("Access Scopes: %s\n", strings.Join(accessScopes, ", "))

	// Create Google authorizer.
	authorizer, err := oauth2.NewAuthorizer(applicationName, accessScopes)
	if err != nil {
		log.Fatalf("Unable to create authorizer: %v", err)
	}

	client, err := authorizer.GetClient()
	if err == nil && client == nil {
		err = xerrors.New("no client returned")
	}
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	for _, scope := range accessScopes {
		if scope == "drive" {
			service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
			if err != nil {
				log.Fatalf("Unable to get drive service object from client: %v", err)
			}

			rows, err := service.Files.List(). /*Fields("id, description").*/ Do()
			if err != nil {
				log.Fatalf("Unable to list files: %v", err)
			}

			fmt.Println("First few drive items:")
			for index, item := range rows.Items {
				if index > 4 {
					break
				}
				fmt.Printf("  %s\n", item.Title)
			}

			break
		}
	}

	fmt.Println("authorize finished")
}
