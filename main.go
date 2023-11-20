package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	// Update with your credentials and URLs
	oauth2Config = &oauth2.Config{
		ClientID:     "go-client",
		ClientSecret: "blhZ1ueOSecw8batVFqHQEo2zts6WwdopXdYTfko62wX19Tu",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"email", "profile", "user", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://xreg-vidm01.vmwpslab.pl/SAAS/auth/oauth2/authorize",
			TokenURL: "https://xreg-vidm01.vmwpslab.pl/SAAS/auth/oauthtoken",
		},
	}
	// Random string for state validation
	oauthStateString = "random"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	fmt.Println("Started running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/login">vIDM Log In</a></body></html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	URL := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, URL, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response := fmt.Sprintf("Access Token: %s", token.AccessToken)
	fmt.Fprintf(w, response)
}
