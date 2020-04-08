package main

import (
	"github.com/yfedoruck/webchat/pkg/env"
	"github.com/yfedoruck/webchat/pkg/fail"
	"github.com/yfedoruck/webchat/pkg/web"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(env.AppPath()+"/templates", t.filename)))
	})
	data := map[string]interface{}{"Host": r.Host}
	if authCookie, err := r.Cookie("auth"); err == nil {
		user := web.Cookie{}
		user.Decode(authCookie.Value)

		data["UserData"] = user
		data["Socket"] = env.Conf().Socket
		// data.user = user
	}

	err := t.templ.Execute(w, data)
	fail.Check(err)
}

func logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.RemoveCookie(w)
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func main() {
	// f, _ := os.Create("/var/log/golang/golang-server.log")
	// defer f.Close()
	// log.SetOutput(f)

	r := newRoom()
	go r.run()

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	http.Handle("/signin", &templateHandler{filename: "login.html"})

	// facebook
	http.HandleFunc("/auth/login/facebook", handleFacebookLogin)
	http.HandleFunc(fbCallback, handleFacebookCallback)

	// google
	http.HandleFunc("/auth/login/google", handleGoogleLogin)
	http.HandleFunc(googleCallback, handleGoogleCallback)

	// github
	http.HandleFunc("/auth/login/github", handleGitHubLogin)
	http.HandleFunc(githubCallback, handleGitHubCallback)

	http.Handle("/logout", logout())

	// start the web server
	log.Println("Starting web server on", env.Port())
	if err := http.ListenAndServe(":"+env.Port(), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
