package handlers

import (
	"encoding/json"
	"fmt"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func handleMessages(client *pkg.Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s Data read error: %v\n", client.Type, err)
			break
		}
		fmt.Printf("Received %s Data: %s\n", client.Type, string(msg))
		client.Timer.Reset(10 * time.Second)
	}
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
		var data models.SystemInfo
		data.Status = 0
		data.TotalSensor = 0
		data.RamConsume = ""
		data.CpuConsume = ""
		data.DHTSensor = 0
		data.BMP180Sensor = 0
		data.RainSensor = 0
		jasonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshal data: ", err)
			return
		}
		fmt.Println("WebSocket inactive for 10 seconds, sending status 0")
		client.Conn.WriteMessage(websocket.TextMessage, jasonData)
	})

	pkg.AddClient(client)
	defer pkg.RemoveClient(client)
	defer client.Conn.Close()
	defer client.Timer.Stop()

	handleMessages(client)
}
