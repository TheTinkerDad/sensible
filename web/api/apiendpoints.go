package api

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/releaseinfo"
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Provides some minimal info about the process
func DoInfo() SimpleApiResult {

	return SimpleApiResult{Result: Body{Result: fmt.Sprintf("Sensible v%s Running...", releaseinfo.Version)}, Status: http.StatusOK}
}

// Shuts down the Sensible server
func DoShutdown(Server *http.Server) SimpleApiResult {

	log.Info("Shutting down...")
	Server.Shutdown(context.TODO())
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}

// Pauses MQTT sensor updates
func DoPauseMqtt() SimpleApiResult {

	log.Info("MQTT Sensor updates paused.")
	mqtt.Paused = true
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}

// Resumes MQTT sensor updates
func DoResumeMqtt() SimpleApiResult {

	log.Info("MQTT Sensor updates resumed.")
	mqtt.Paused = false
	return SimpleApiResult{Result: Body{Result: "OK"}, Status: http.StatusOK}
}
