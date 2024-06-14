package handlers

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (h *handlers) HandleConnectionWs(c *gin.Context) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &pkg.Client{Conn: conn}
	client.Timer = time.AfterFunc(10*time.Second, func() {
		fmt.Println("WebSocket inactive for 10 seconds, sending status 0")
		client.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status": 0}`))
	})

	pkg.AddClient(conn)
	defer pkg.RemoveClient(conn)
	defer client.Conn.Close()
	defer client.Timer.Stop()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		client.Timer.Reset(10 * time.Second)
		pkg.Broadcast <- msg
	}
}
