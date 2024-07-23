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
	Clients           = make(map[*websocket.Conn]bool)
	SensorClients     = make(map[*Client]bool)
	SystemClients     = make(map[*Client]bool)
	ClientsByUsername = make(map[string]map[*Client]bool)
	Broadcast         = make(chan []byte)
	mutex             = &sync.Mutex{}
)

type Client struct {
	Conn     *websocket.Conn
	Timer    *time.Timer
	Type     string
	Username string
}

func AddClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	Clients[client.Conn] = true
	if client.Type == "sensor" {
		SensorClients[client] = true
	} else if client.Type == "system" {
		SystemClients[client] = true
	}

	if _, ok := ClientsByUsername[client.Username]; !ok {
		ClientsByUsername[client.Username] = make(map[*Client]bool)
	}
	ClientsByUsername[client.Username][client] = true
}

func RemoveClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(Clients, client.Conn)
	if client.Type == "sensor" {
		delete(SensorClients, client)
	} else if client.Type == "system" {
		delete(SystemClients, client)
	}

	if clients, ok := ClientsByUsername[client.Username]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(ClientsByUsername, client.Username)
		}
	}
	client.Conn.Close()
}

func BroadcastToSensors(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range SensorClients {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func BroadcastToSystems(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range SystemClients {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func BroadcastToUsername(username string, message []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	if clients, ok := ClientsByUsername[username]; ok {
		for client := range clients {
			client.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
