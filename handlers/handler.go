package handlers

import (
	"encoding/json"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlers struct {
	services service.Service
}

func NewHandler(services service.Service) *handlers {
	return &handlers{
		services: services,
	}
}

func (h *handlers) CreateData(c *gin.Context) {
	if c.Request.Body != nil || c.Request.ContentLength > 0 {
		var dataAntares models.SensorData
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&dataAntares); err != nil {
			c.JSON(400, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Bad Request",
			})
			return
		}
		err := h.services.InsertData(c, dataAntares)
		if err != nil {
			c.JSON(500, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Internal Server Error",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": http.StatusOK,
			"msg":  "Success",
		})
		return
	} else {
		c.JSON(400, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Request Body is Empty",
		})
	}
}
