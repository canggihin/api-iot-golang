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
	defer func() {
		if client.Timer != nil {
			client.Timer.Stop()
		}
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s Data read error: %v\n", client.Type, err)
			return // Exit the loop and let defer clean up
		}
		fmt.Printf("Received %s Data: %s\n", client.Type, string(msg))

		// Reset the timer on message receipt
		if client.Timer != nil {
			client.Timer.Reset(60 * time.Second)
		}
	}
}

func setupTimer(client *pkg.Client) {
	client.Timer = time.AfterFunc(60*time.Second, func() {
		sendStatus(client)
		client.Timer.Reset(60 * time.Second) // Immediately reset the timer
	})
}

func sendStatus(client *pkg.Client) {
	var data models.SystemInfo
	data.Status = 0
	data.TotalSensor = 0
	data.RamConsume = ""
	data.CpuConsume = ""
	data.DHTSensor = 0
	data.BMP180Sensor = 0
	data.RainSensor = 0
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data: ", err)
		return
	}
	fmt.Println("WebSocket inactive for 10 seconds, sending status 0")
	client.Conn.WriteMessage(websocket.TextMessage, jsonData)
}

func handleCWs(c *gin.Context, clientType string) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &pkg.Client{Conn: conn, Type: clientType}
	setupTimer(client)

	pkg.AddClient(client)
	defer pkg.RemoveClient(client)
	defer client.Conn.Close()

	handleMessages(client)
}

func (h *handlers) HandleWsSensor(c *gin.Context) {
	handleCWs(c, "sensor")
}

func (h *handlers) HandleWsSystem(c *gin.Context) {
	handleCWs(c, "system")
}
