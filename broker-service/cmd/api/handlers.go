package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	c.Header("Contect-Type", "application/json")
	c.JSON(http.StatusAccepted, out)
}
