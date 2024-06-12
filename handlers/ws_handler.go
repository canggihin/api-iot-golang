package handlers

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/pkg"

	"github.com/gin-gonic/gin"
)

func (h *handlers) HandleConnectionWs(c *gin.Context) {
	conn, err := pkg.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	pkg.Clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(pkg.Clients, conn)
			return
		}
		pkg.Broadcast <- msg
	}
}
