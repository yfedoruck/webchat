package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
)

func basePath() string {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}
	return filepath.Dir(b)
}
func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type tplData struct {
	user
	Host string
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(basePath()+"/templates", t.filename)))
	})
	data := map[string]interface{}{"Host": r.Host}
	if authCookie, err := r.Cookie("auth"); err == nil {
		user := cookie{}
		user.decode(authCookie.Value)

		// dec := json.NewDecoder(strings.NewReader(authCookie.Value))
		// err := dec.Decode(&user)
		data["UserData"] = user
		// data.user = user
		fmt.Println(data)
	}
	// fmt.Println(data)

	err := t.templ.Execute(w, data)
	check(err)
}

func logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()
	r := newRoom()

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	// http.HandleFunc("/", handleMain)
	go r.run()
	http.Handle("/room", r)

	http.HandleFunc("/signin", handleMain)
	http.HandleFunc("/login", handleFacebookLogin)
	http.Handle("/logout", logout())
	http.HandleFunc("/oauth2callback", handleFacebookCallback)

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServeTLS(*addr, filepath.Join(basePath()+"/server.rsa.crt"), filepath.Join(basePath()+"/server.rsa.key"), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	// if err := http.ListenAndServe(*addr, nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }
}
