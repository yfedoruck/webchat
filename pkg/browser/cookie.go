package browser

import (
	"encoding/base64"
	"encoding/json"
	"github.com/yfedoruck/webchat/pkg/fail"
	"net/http"
)

type Cookie struct {
	Name      string
	AvatarURL string
}

func (c Cookie) encode() string {
	js, err := json.Marshal(c)
	fail.Check(err)

	return base64.StdEncoding.EncodeToString(js)
}

func (c *Cookie) Decode(arg string) {
	js, err := base64.StdEncoding.DecodeString(arg)
	fail.Check(err)

	err = json.Unmarshal(js, &c)
	fail.Check(err)
}

func (c Cookie) Set(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: c.encode(),
		Path:  "/",
	})
}

func RemoveCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
