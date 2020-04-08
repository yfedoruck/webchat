package chat

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *Room
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		// _, msg, err := c.socket.ReadMessage()
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL = c.userData["avatar_url"].(string)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		// err := c.socket.WriteMessage(websocket.TextMessage, msg)
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
