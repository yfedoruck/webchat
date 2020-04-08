package main

import (
	"github.com/yfedoruck/webchat/pkg/chat"
	"github.com/yfedoruck/webchat/pkg/env"
	"github.com/yfedoruck/webchat/pkg/web"
	"log"
	"net/http"
)

func main() {
	// f, _ := os.Create("/var/log/golang/golang-server.log")
	// defer f.Close()
	// log.SetOutput(f)

	room := chat.NewRoom()
	go room.Run()

	http.Handle("/", web.MustAuth(&web.TemplateHandler{Filename: "chat.html"}))
	http.Handle("/room", room)

	http.Handle("/signin", &web.TemplateHandler{Filename: "login.html"})

	// facebook
	http.HandleFunc("/auth/login/facebook", web.HandleFacebookLogin)
	http.HandleFunc(web.FbCallback, web.HandleFacebookCallback)

	// google
	http.HandleFunc("/auth/login/google", web.HandleGoogleLogin)
	http.HandleFunc(web.GoogleCallback, web.HandleGoogleCallback)

	// github
	http.HandleFunc("/auth/login/github", web.HandleGitHubLogin)
	http.HandleFunc(web.GithubCallback, web.HandleGitHubCallback)

	http.Handle("/logout", web.Logout())

	// start the web server
	log.Println("Starting web server on", env.Port())
	if err := http.ListenAndServe(":"+env.Port(), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
