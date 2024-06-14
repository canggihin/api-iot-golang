package pkg

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan []byte)
	mutex     = &sync.Mutex{}
)

type Client struct {
	Conn  *websocket.Conn
	Timer *time.Timer
}

func AddClient(conn *websocket.Conn) {
	mutex.Lock()
	Clients[conn] = true
	mutex.Unlock()
}

func RemoveClient(conn *websocket.Conn) {
	mutex.Lock()
	delete(Clients, conn)
	mutex.Unlock()
}
