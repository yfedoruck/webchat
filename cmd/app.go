package main

import (
	"github.com/yfedoruck/webchat/pkg/chat"
	"github.com/yfedoruck/webchat/pkg/web"
)

type App struct {
	server *web.Server
	room   *chat.Room
}

func (a App) Init() {
	a.room = chat.NewRoom()
	a.server = web.NewServer(a.room)
}

func (a App) Run() {
	go a.room.Run()
	a.server.Start()
}
