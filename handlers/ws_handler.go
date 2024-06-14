package handlers

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func handleMessages(client *pkg.Client, clientType string) {
	_, msg, err := client.Conn.ReadMessage()
	if err != nil {
		fmt.Printf("%s Data read error: %v\n", clientType, err)
		return
	}
	fmt.Printf("Received %s Data: %s\n", clientType, string(msg))
}

func (h *handlers) HandleWsSensor(c *gin.Context) {
	handleCWs(c, "sensor")
}

func (h *handlers) HandleWsSystem(c *gin.Context) {
	handleCWs(c, "system")
}

func handleCWs(c *gin.Context, clientType string) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &pkg.Client{Conn: conn, Type: clientType}
	client.Timer = time.AfterFunc(10*time.Second, func() {
		fmt.Println("WebSocket inactive for 10 seconds, sending status 0")
		client.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status": 0}`))
	})

	pkg.AddClient(client)
	defer pkg.RemoveClient(client)
	defer client.Conn.Close()
	defer client.Timer.Stop()

	for {
		client.Timer.Reset(10 * time.Second)
		handleMessages(client, clientType)
	}
}
