package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/pkg"
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

		if dataAntares.Temperature == 0 || dataAntares.Humidity == 0 || dataAntares.Pressure == 0 {
			c.JSON(400, gin.H{
				"code": http.StatusBadRequest,
				"msg":  fmt.Sprintf("Invalid sensor data: %+v", dataAntares),
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

func (h *handlers) GetData(c *gin.Context) {
	data, err := h.services.GetData(c)
	if err != nil {
		if err.Error() == "data not found" {
			c.JSON(404, gin.H{
				"code": http.StatusNotFound,
				"msg":  "Data Not Found",
			})
			return
		}
		log.Println(err)
		c.JSON(500, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
		"data": data,
	})
}

func (h *handlers) GetDataByDay(c *gin.Context) {
	data, err := h.services.GetDataByDay(c)
	if err != nil {
		if err.Error() == "data not found" {
			c.JSON(404, gin.H{
				"code": http.StatusNotFound,
				"msg":  "Data Not Found",
			})
			return
		}
		log.Println(err)
		c.JSON(500, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
		"data": data,
	})
}

func (h *handlers) SuccessConnectedDevice(c *gin.Context) {
	pkg.BroadcastToSystems([]byte("Connected Device"))
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
	})
}

func (h *handlers) GetSystemInfo(c *gin.Context) {
	var dataReq models.SystemInfo
	if err := c.ShouldBindJSON(&dataReq); err != nil {
		c.JSON(400, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Bad Request",
		})
		return

	}

	data, err := h.services.ProsesMessage(context.Background(), dataReq)
	if err != nil {
		if err.Error() == "data not found" {
			c.JSON(404, gin.H{
				"code": http.StatusNotFound,
				"msg":  "Data Not Found",
			})
			return
		}
		log.Println(err)
		c.JSON(500, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
		"data": data,
	})
}
