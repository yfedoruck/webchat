package main

import (
	"flag"
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
		t.templ = template.Must(template.ParseFiles(filepath.Join(basePath()+"/templates", t.filename)))
	})
	data := map[string]interface{}{"Host": r.Host}
	if authCookie, err := r.Cookie("auth"); err == nil {
		user := cookie{}
		user.decode(authCookie.Value)

		data["UserData"] = user
		// data.user = user
	}

	err := t.templ.Execute(w, data)
	check(err)
}

func logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		removeCookie(w)
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()
	r := newRoom()
	go r.run()

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	http.Handle("/signin", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/login/facebook", handleFacebookLogin)
	http.HandleFunc("/oauth2callback", handleFacebookCallback)

	http.HandleFunc("/auth/login/google", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)

	http.Handle("/logout", logout())

	// start the web server
	log.Println("Starting web server on", *addr)
	// if err := http.ListenAndServe(*addr, nil); err != nil {
	if err := http.ListenAndServeTLS(*addr, filepath.Join(basePath()+"/server.rsa.crt"), filepath.Join(basePath()+"/server.rsa.key"), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
