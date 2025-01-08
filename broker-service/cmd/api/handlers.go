package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPlayload `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPlayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	c.IndentedJSON(http.StatusAccepted, payload)
}

func (app *Config) HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	if err := c.BindJSON(&requestPayload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(c, requestPayload.Auth)
	case "log":
		app.logItem(c, requestPayload.Log)
	default:
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
	}
}

func (app *Config) authenticate(c *gin.Context, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		log.Println(err)
	}

	// call the service
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Contect-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	// make sure we get back the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		c.IndentedJSON(res.StatusCode, gin.H{"error": "invalid credentials"})
		return
	} else if res.StatusCode != http.StatusAccepted {
		c.IndentedJSON(res.StatusCode, gin.H{"error": "error calling auth service"})
		return
	}

	// create a variable we'll read res.Body into
	var jsonFromService jsonResponse

	if err = json.NewDecoder(res.Body).Decode(&jsonFromService); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if jsonFromService.Error {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated !"
	payload.Data = jsonFromService.Data

	c.IndentedJSON(http.StatusAccepted, payload)
}

func (app *Config) logItem(c *gin.Context, entry LogPlayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		c.IndentedJSON(response.StatusCode, gin.H{"error": "not accepted"})
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	c.IndentedJSON(http.StatusAccepted, payload)
}
