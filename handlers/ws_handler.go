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

// handleMessages function to handle incoming WebSocket messages
func handleMessages(client *pkg.Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s Data read error: %v\n", client.Type, err)
			return // Exit the loop and let defer clean up
		}
		fmt.Printf("Received %s Data: %s\n", client.Type, string(msg))

		// Reset the timer on message receipt
		client.Timer.Reset(60 * time.Second)
	}
}

// setupTimer function to set up a timer for the client
func setupTimer(client *pkg.Client) {
	client.Timer = time.AfterFunc(60*time.Second, func() {
		sendStatus(client)
	})
}

// sendStatus function to send a status update if no message is received within the timeout period
func sendStatus(client *pkg.Client) {
	var data models.SystemInfo
	data.BatteryLevel = 0
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
	fmt.Println("WebSocket inactive for 60 seconds, sending status 0")
	client.Conn.WriteMessage(websocket.TextMessage, jsonData)

	// Reset the timer to check for inactivity again after sending the status
	client.Timer.Reset(60 * time.Second)
}

// handleCWs function to handle WebSocket connections
func handleCWs(c *gin.Context, clientType string) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	username := c.Param("username")
	client := &pkg.Client{Conn: conn, Type: clientType, Username: username}
	setupTimer(client)

	pkg.AddClient(client)
	defer pkg.RemoveClient(client)
	defer client.Conn.Close()

	handleMessages(client)
}

// HandleWsSensor function to handle sensor WebSocket connections
func (h *handlers) HandleWsSensor(c *gin.Context) {
	handleCWs(c, "sensor")
}

// HandleWsSystem function to handle system WebSocket connections
func (h *handlers) HandleWsSystem(c *gin.Context) {
	handleCWs(c, "system")
}
