package web

import (
	"encoding/json"
	"fmt"
	"github.com/yfedoruck/webchat/pkg/browser"
	"github.com/yfedoruck/webchat/pkg/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

const GoogleCallback string = "/callback"

var (
	googleConf  *oauth2.Config
	googleState = "pseudorandom"
)

func init() {
	c := Config{}
	c.Set("google")
	googleConf = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  env.Conf().Host + GoogleCallback,
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

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleConf.AuthCodeURL(googleState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// ***** //
	googleUser := googleUser{}
	_ = json.Unmarshal(content, &googleUser)

	browser.Cookie{
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	}.Set(w)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func getUserInfo(state string, code string) ([]byte, error) {
	if state != googleState {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleConf.Exchange(oauth2.NoContext, code)
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
