package handlers

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/pkg"

	"github.com/gin-gonic/gin"
)

func (h *handlers) HandleConnectionWs(c *gin.Context) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	pkg.AddClient(conn)
	defer pkg.RemoveClient(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		pkg.Broadcast <- msg
	}
}
