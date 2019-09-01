# authorize

The `authorize` command manually authorizes an application.

This functionality is not strictly required as the `oauth.authorizer` object
will automagically invoke the authorization protocol as necessary.
This is a good way to verify the authorization secret for an application.

## Usage

```shell script
authorize <applicationName> <accessScope>*
```

The `applicationName` must be the name used when configuring the keys using
the Google developer console
(see [documentation for the `oauth2` package](../../oauth2/README.md)).

An `accessScope` represents a permission that a user must agree to
during the authorization process.
Normally [access scopes](https://developers.google.com/identity/protocols/googlescopes)
are represented as URLs.
When invoking `authorize` just use the final path element.
For example, represent access scope `https://www.googleapis.com/auth/drive` as `drive`.

* If no access scopes are provided the single `drive` scope is used.
* If `drive` is available `authorize` will attempt to list
the first five items in the top-level directory as a kind of "ping".
