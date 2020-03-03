package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	c := config{}
	c.set("google")
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  c.RedirectURL,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type googleUser struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

var (
	// TODO: randomize it
	oauthStateStringGoogle = "pseudo-random"
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateStringGoogle)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// ***** //
	googleUser := googleUser{}
	_ = json.Unmarshal(content, &googleUser)

	cookie{
		googleUser.Name,
		googleUser.Picture,
	}.set(w)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateStringGoogle {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
