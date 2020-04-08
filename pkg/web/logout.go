package web

import (
	"github.com/yfedoruck/webchat/pkg/browser"
	"net/http"
)

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		browser.RemoveCookie(w)
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
