# OAuth2

Handle authentication for Google API.

## Project Setup

*The following are overly simplistic instructions for a complex setup.
Please refer to the appropriate Google documentation.*

The [Google Deverlopers Console](https://console.developers.google.com/apis/dashboard)
is used to configure access to Google APIs.
*You probably need to be a domain owner to do this.*

Use the developer console to create a new project if necessary
(there appears to be a `QuickStart` project by default).
Once you have a project, set the developer console to that project using the drop-down.

Use **ENABLE APIS AND SERVICES** to select Google APIs that are required for the project.
An API will include multiple *Access Scopes*.

Then use the **OAuth consent screen** to configure:

* **Application type**: Internal
* **Application name**: *whatever is appropriate*
* Use the **Add scope** button to add whatever scopes are required for your
specific project.

There are a lot of different Google APIs.
These are specified as [*Access Scopes*](https://developers.google.com/identity/protocols/googlescopes).

## Credentials

OAuth2 requires credentials from Google specific to your application.
The following links:

* [
Using OAuth 2.0 to Access Google APIs](https://developers.google.com/identity/protocols/OAuth2)
* [
OAuth 2.0 for Mobile & Desktop Apps](https://developers.google.com/identity/protocols/OAuth2InstalledApp#identify-access-scopes)

provide background and instructions for acquiring and using credentials.
The code in this package implements **Option 2: Loopback IP address**.

With your project selected in the developer console use the **Create credentials**
button start the credential creation process.
Configure the new credentials:

* Choose **OAuth client id**
* **Application type**: Other (to implement **Loopback IP address**)
* **Name**: *name for the credential as there can be multiple*

The new credentials should show up on the developer console.

### Credential data

Once you have the credential you can select it to see the **Client ID** and **Client secret**.
Credential information must be available to your application.

This package loads credential information from a JSON file stored in `~/.ssh`.
You could also embed the information into your program.
Acquire the JSON file via the **DOWNLOAD JSON** button from the credential information page.
Rename the file to `<project>-secret.json` and move it to `~/.ssh`.

## Usage

There are three components for OAuth2 usage for Google APIs:

* the `authorize` application used to acquire an access token and
* the `oauth2` package which creates a client.

There are other package(s) in this repository that make working with
various Google APIs somewhat more straightforward (YMMV).

### `authorize`

The `authorize` application has two purposes:

* specify a URL for the user to use for authorization.
This URL, constructed using credential data, is served by Google and
technically it provides both authentication and authorization and
* provide a local HTTP server to receive the access token from Google.

### `oauth2`