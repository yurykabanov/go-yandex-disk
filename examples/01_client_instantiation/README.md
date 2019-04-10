# Yandex.Disk Client instantiation

## How to use

There're three ways to create new client:
- using your own properly configured `http.Client`
- using `oauth2.Config` and `oauth2.Token`
- using just an access token

### Using `http.Client`

```
import (
	"context"
	"net/http"
	
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"

	"github.com/yurykabanov/go-yandex-disk"
)

// ...

oauthConfig := &oauth2.Config{
    ClientID:     "CLIENT-ID",
    ClientSecret: "CLIENT-SECRET",
    RedirectURL:  "http://your-app.tld/auth/yandex/callback",
    Endpoint:     yandex.Endpoint,
}

// Retrieve token from whatever storage
token := oauth2.Token{
    TokenType:    "OAuth",
    AccessToken:  "ACCESS-TOKEN",
    RefreshToken: "REFRESH-TOKEN",
    Expiry:       /* proper time.Time */,
}

source := oauthConfig.TokenSource(context.TODO(), &token)

// You can configure http.Client as needed
httpClient := &http.Client{
    Transport: &oauth2.Transport{
        Source: source,
    },
}

client := yadisk.New(httpClient)
```

Note the following:
1. You **MUST** use `oauth2.Transport` with proper token source.
2. The token **MUST** have `TokenType` equals to `OAuth` (instead of
default `Bearer`). It is descibed in
[Yandex.Disk API documentation](https://tech.yandex.com/disk/api/concepts/quickstart-docpage/).
3. You **SHOULD** replace `context.TODO()` in token source with more
appropriate one.

### Using `oauth2.Config` and `oauth2.Token`

It is just the same as previous one, just a little bit shorter:

```
oauthConfig := &oauth2.Config{
    /* ... */
}

// Retrieve token from whatever storage
token := &oauth2.Token{
    TokenType:    "OAuth",
    AccessToken:  "ACCESS-TOKEN",
    RefreshToken: "REFRESH-TOKEN",
    Expiry:       /* proper time.Time */,    
}

client := yadisk.NewFromConfigAndToken(oauthConfig, token, context.TODO())
```

You should note exactly the same things as in previous method.

### Using just an access token

The most simplest way:
```
client := yadisk.NewFromAccessToken(*accessToken)
```

But it **WILL NOT** refresh your token automatically as two previous methods.

## Building
Build example as following:
```bash
go build -o ./bin/yandex-client ./examples/01_client_instantiation/main.go

```

## Running
It doesn't really perform any useful action, but you can run it:
```bash
./bin/yandex-client                               \ 
    -using=own-http-client                        \
    -client-id=YOUR-APPLICATION-CLIENT-ID         \
    -client-secret=YOUR-APPLICATION-CLIENT-SECRET \
    -access-token=YOUR-ACCESS-TOKEN               \
    -refresh-token=YOUR-REFRESH-TOKEN             \
    -expiry-time="2020-01-02T12:34:56Z"
    
./bin/yandex-client                               \ 
    -using=oauth-config-and-token                 \
    -client-id=YOUR-APPLICATION-CLIENT-ID         \
    -client-secret=YOUR-APPLICATION-CLIENT-SECRET \
    -access-token=YOUR-ACCESS-TOKEN               \
    -refresh-token=YOUR-REFRESH-TOKEN             \
    -expiry-time="2020-01-02T12:34:56Z"
    
./bin/yandex-client                               \ 
    -using=access-token-only                      \
    -access-token=YOUR-ACCESS-TOKEN                            
```
