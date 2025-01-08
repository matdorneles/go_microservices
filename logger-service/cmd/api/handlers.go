package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matdorneles/go_microservices/logger-service/data"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (app *Config) WriteLog(c *gin.Context) {
	// rest json into var
	var requestPayload jsonPayload
	if err := c.BindJSON(&requestPayload); err != nil {
		log.Println(err)
		return
	}

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	if err := app.Models.LogEntry.Insert(event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	c.JSON(http.StatusAccepted, res)
}
