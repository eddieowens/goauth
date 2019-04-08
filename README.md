# GoAuth
Lightweight dependency for local OAuth. 

## What does it do?
Functionally equivalent to Google Cloud SDK's command `gcloud auth login`.

Opens up a user's default browser to authorize your app and returns the token. Meant
to be used for storing OAuth tokens locally on a user's machine.

## Installation
```bash
go get github.com/eddieowens/goauth
```

## Basic Usage
### Synchronously
```go
package main

import (
    "fmt"
    "github.com/eddieowens/goauth"
)

func main() {
    oauth := goauth.OAuth{
        ClientSecret: "some_secret",
        ClientId:     "some_id",
        Scope:        "https://www.googleapis.com/auth/userinfo.email",
    }
    
    token := goauth.Google(oauth).Auth()
    if token.Error != nil {
        panic(token.Error)
    }
    fmt.Println("Access token is " + token.AccessToken)
}
```
The above will block until the user authorizes your app. When `Auth()` is called, a
window in the user's default browser prompting them to select a Google account will
open. Upon selecting an account, the OAuth 2 protocol is followed and a 
`GoogleOAuthToken` struct is returned with the credentials.
### Asynchronously
```go
package main

import (
    "fmt"
    "github.com/eddieowens/goauth"
)

func main() {
    oauth := goauth.OAuth{
        ClientSecret: "some_secret",
        ClientId:     "some_id",
        Scope:        "https://www.googleapis.com/auth/userinfo.email",
    }
	
    err := goauth.Google(oauth).
        AsyncAuth(func(token goauth.GoogleOAuthToken) {
            if token.Error != nil {
                panic(token.Error)
            }
            fmt.Println("Access token is " + token.AccessToken)
        })
    if err != nil {
        panic(err)
    }
}
```
Same as the [synchronous](#synchronously) example, only the entire user flow is
non-blocking.

## Supported OAuth providers
* [Google](https://developers.google.com/identity/protocols/OAuth2)

If you want a specific provider supported, just create an [issue](https://github.com/eddieowens/goauth/issues/new).