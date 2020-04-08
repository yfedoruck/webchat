package web

import "net/http"

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RemoveCookie(w)
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
