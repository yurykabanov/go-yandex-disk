package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

var (
	listenAddr   = flag.String("listen", "localhost:8000", "Listen address")
	clientID     = flag.String("client-id", "", "OAuth Client ID")
	clientSecret = flag.String("client-secret", "", "OAuth Client Secret")
	callbackUrl  = flag.String("callback-url", "http://localhost:8000/auth/yandex/callback", "OAuth Callback URL")
)

type YandexAuth struct {
	oauthConfig *oauth2.Config
}

func NewYandexAuth(clientId, clientSecret, callbackUrl string) *YandexAuth {
	return &YandexAuth{
		oauthConfig: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  callbackUrl,
			Endpoint:     yandex.Endpoint,
		},
	}
}

func main() {
	flag.Parse()

	auth := NewYandexAuth(*clientID, *clientSecret, *callbackUrl)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/yandex/login", auth.Login)
	mux.HandleFunc("/auth/yandex/callback", auth.Callback)

	server := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (h *YandexAuth) generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)

	return state
}

func (h *YandexAuth) Login(w http.ResponseWriter, r *http.Request) {
	state := h.generateStateOauthCookie(w)
	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *YandexAuth) Callback(w http.ResponseWriter, r *http.Request) {
	if state, err := r.Cookie("oauthstate"); err != nil || r.FormValue("state") != state.Value {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad state"))
		return
	}

	code := r.FormValue("code")

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad code"))
		return
	}

	w.WriteHeader(http.StatusOK)

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}

	w.Write(data)
}
