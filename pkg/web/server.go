package web

import (
	"github.com/yfedoruck/webchat/pkg/chat"
	"github.com/yfedoruck/webchat/pkg/env"
	"log"
	"net/http"
)

type Server struct {
	Port string
	room *chat.Room
}

func NewServer(room *chat.Room) *Server {
	s := &Server{}
	s.Port = env.Port()
	s.room = room
	s.routers()
	return s
}

func (s *Server) Start() {
	log.Println("Starting web server on", s.Port)
	if err := http.ListenAndServe(":"+s.Port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func (s *Server) routers() {
	// f, _ := os.Create("/var/log/golang/golang-server.log")
	// defer f.Close()
	// log.SetOutput(f)

	http.Handle("/", MustAuth(&TemplateHandler{Filename: "chat.html"}))
	http.Handle("/room", s.room)

	http.Handle("/signin", &TemplateHandler{Filename: "login.html"})

	// facebook
	http.HandleFunc("/auth/login/facebook", HandleFacebookLogin)
	http.HandleFunc(FbCallback, HandleFacebookCallback)

	// google
	http.HandleFunc("/auth/login/google", HandleGoogleLogin)
	http.HandleFunc(GoogleCallback, HandleGoogleCallback)

	// github
	http.HandleFunc("/auth/login/github", HandleGitHubLogin)
	http.HandleFunc(GithubCallback, HandleGitHubCallback)

	http.Handle("/logout", Logout())
}
