package pkg

import (
	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan []byte)
)
