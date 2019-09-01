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
It is likely that you will always need access to the `drive` API.

Then use the **OAuth consent screen** to configure:

* **Application type**: `Internal`
* **Application name**: *name of your application*
* Use the **Add scope** button to add whatever scopes are required for your
specific project, out of the set enabled earlier.

There are
[a lot of different Google APIs](https://developers.google.com/identity/protocols/googlescopes).

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
button to start the credential creation process.
Configure the new credentials:

* Choose **OAuth client id**
* **Application type**: `Other` (to implement **Loopback IP address**)
* **Name**: *name for the credential* (there can be multiple)

The new credentials should show up on the developer console.

### Credential data

Once you have the credential you can select it to see the **Client ID** and **Client secret**.
Credential information must be available to your application.

This package loads credential information from a JSON file stored in `~/.ssh`.
You could also embed the information into your program,
but this package does not support that option.
Acquire the JSON file via the **DOWNLOAD JSON** button from the credential information page.
Rename the file to `<applicationName>-secret.json` and move it to `~/.ssh`.

## Usage

### `authorize`

The `authorize` application, provided in the command section of this repository,
can be used to manually acquire an authorization token for an application.
In general this is not required but it provides a good test of basic authorization functionality.

### `oauth2`

This package provides basic authorization for a package.
When properly configured it will look for an existing access token in `~/.ssh/<applicationName>-token.json`.
If not found a browser window will (hopefully) be initiated for user interaction
and a temporary loopback HTTP server will be called from the browser with the access token data.
The token file will be created in `~/.ssh` for current and later use.

## Caveats

Using `plain` mode for code verifier challenge.
Tried to implement SHA256 version but was unable to get it to work.
