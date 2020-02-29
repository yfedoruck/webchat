package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "245616582805552",
		ClientSecret: "b6568aba9223f933ba52bcbf77aa773d",
		RedirectURL:  "https://localhost:8080/oauth2callback",
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

type user struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

const htmlIndex = `<html><body>
Logged in with <a href="/login">facebook</a>
</body></html>
`

func getUser(byt []byte) user {
	// dat := make(map[string]interface{})
	user := user{}

	if err := json.Unmarshal(byt, &user); err != nil {
		log.Panic(err)
	}

	fmt.Println(user.Id, user.Name)

	return user
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	// url := Url.String()
	authUrl := oauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
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

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// user := getUser(response)
	// authCookie, err := json.Marshal(user)
	authCookie := base64.StdEncoding.EncodeToString(response)
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookie,
		Path:  "/",
	})

	log.Printf("parseResponseBody: %s\n", string(response))
	res := fmt.Sprintf("%s", response)
	fmt.Println(res)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main2() {
	// http.HandleFunc("/", handleMain)
	// http.HandleFunc("/login", handleFacebookLogin)
	// http.HandleFunc("/oauth2callback", handleFacebookCallback)
	// fmt.Print("Started running on http://localhost:9090\n")
	// log.Fatal(http.ListenAndServe(":9090", nil))
}
