package api

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/releaseinfo"
	"context"
	"fmt"
	"log"
	"net/http"
)

// Provides some minimal info about the process
func DoInfo() SimpleApiResult {

	return SimpleApiResult{Result: Body{Result: fmt.Sprintf("Sensible v%s Running...", releaseinfo.Version)}, Status: http.StatusOK}
}

// Shuts down the Sensible server
func DoShutdown(Server *http.Server) SimpleApiResult {

	log.Println("Shutting down...")
	Server.Shutdown(context.TODO())
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}

// Pauses MQTT sensor updates
func DoPauseMqtt() SimpleApiResult {

	log.Println("MQTT Sensor updates paused.")
	mqtt.Paused = true
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}

// Resumes MQTT sensor updates
func DoResumeMqtt() SimpleApiResult {

	log.Println("MQTT Sensor updates resumed.")
	mqtt.Paused = false
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}
