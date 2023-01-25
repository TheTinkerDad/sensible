package api

import (
	"TheTinkerDad/sensible/mqtt"
	"context"
	"log"
	"net/http"
)

// SimpleApiResult Holds data about the result of a simple operation
type SimpleApiResult struct {
	Result string
}

// Provides some minimal info about the process
func Info() SimpleApiResult {

	return SimpleApiResult{Result: "Sensible v0.1 Running..."}
}

// Shuts down the Sensible server
func Shutdown(Server *http.Server) SimpleApiResult {

	log.Println("Shutting down...")
	Server.Shutdown(context.TODO())
	return SimpleApiResult{Result: "OK"}
}

// Pauses MQTT sensor updates
func PauseMqtt() SimpleApiResult {

	log.Println("MQTT Sensor updates paused.")
	mqtt.Paused = true
	return SimpleApiResult{Result: "OK"}
}

// Resumes MQTT sensor updates
func ResumeMqtt() SimpleApiResult {

	log.Println("MQTT Sensor updates resumed.")
	mqtt.Paused = false
	return SimpleApiResult{Result: "OK"}
}
