package main

import (
	"fmt"
	"github.com/yfedoruck/webchat/pkg/env"
	"github.com/yfedoruck/webchat/pkg/web"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

const githubCallback string = "/oauth2github"

var (
	oauthConf *oauth2.Config
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
)

func init() {
	c := web.Config{}
	c.Set("github")
	oauthConf = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  env.Conf().Host + githubCallback,
		Scopes:       []string{"user:username,avatar_url"},
		Endpoint:     githuboauth.Endpoint,
	}

	fmt.Println(oauthConf)
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//
	web.Cookie{
		Name:      *user.Login,
		AvatarURL: *user.AvatarURL,
	}.Set(w)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
