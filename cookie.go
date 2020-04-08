package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/yfedoruck/webchat/pkg/fail"
	"net/http"
)

type cookie struct {
	Name      string
	AvatarURL string
}

func (c cookie) encode() string {
	js, err := json.Marshal(c)
	fail.Check(err)

	return base64.StdEncoding.EncodeToString(js)
}

func (c *cookie) decode(arg string) {
	js, err := base64.StdEncoding.DecodeString(arg)
	fail.Check(err)

	err = json.Unmarshal(js, &c)
	fail.Check(err)
}

func (c cookie) set(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: c.encode(),
		Path:  "/",
	})
}

func removeCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
