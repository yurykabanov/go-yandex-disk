package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"

	"github.com/yurykabanov/go-yandex-disk"
)

var (
	using = flag.String("using", "", "Which example to use")

	clientID     = flag.String("client-id", "", "OAuth Client ID")
	clientSecret = flag.String("client-secret", "", "OAuth Client Secret")

	accessToken  = flag.String("access-token", "", "Access Token")
	refreshToken = flag.String("refresh-token", "", "Refresh Token")
	expiryTime   = flag.String("expiry-time", "", "Expiry Time in RFC3339 format")
)

func main() {
	flag.Parse()

	var client *yadisk.Client

	switch *using {
	case "own-http-client":
		client = UsingYourOwnHttpClient()
	case "oauth-config-and-token":
		client = UsingOauthConfigAndToken()
	case "access-token-only":
		client = UsingAccessToken()
	default:
		fmt.Println("You should specify -using=<...> with either 'own-http-client', 'access-token-only' or 'oauth-config-and-token'")
		os.Exit(1)
	}

	// Do some stuff with client
	_ = client
}

// 1. Configure/reuse your own HTTP Client. Make sure to specify OAuth2 Transport.
func UsingYourOwnHttpClient() *yadisk.Client {
	oauthConfig := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Endpoint:     yandex.Endpoint,
	}

	expiry, err := time.Parse(time.RFC3339, *expiryTime)
	if err != nil {
		log.Fatal("invalid expiry time")
	}

	// Retrieve token from whatever storage
	token := oauth2.Token{
		TokenType:    "OAuth", // NOTE: it's important to use "OAuth" here instead of "Bearer"
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		Expiry:       expiry,
	}

	source := oauthConfig.TokenSource(context.TODO(), &token)

	// You can configure http.Client as needed
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: source,
		},
	}

	return yadisk.New(httpClient)
}

// 2. Small shortcut for your way #1
func UsingOauthConfigAndToken() *yadisk.Client {
	oauthConfig := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Endpoint:     yandex.Endpoint,
	}

	expiry, err := time.Parse(time.RFC3339, *expiryTime)
	if err != nil {
		log.Fatal("invalid expiry time")
	}

	// Retrieve token from whatever storage
	token := &oauth2.Token{
		TokenType:    "OAuth", // NOTE: it's important to use "OAuth" here instead of "Bearer"
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		Expiry:       expiry,
	}

	return yadisk.NewFromConfigAndToken(oauthConfig, token, context.TODO())
}

// 3. Short and simple way to use client with only an access token. It will not be refreshed though.
func UsingAccessToken() *yadisk.Client {
	return yadisk.NewFromAccessToken(*accessToken)
}
