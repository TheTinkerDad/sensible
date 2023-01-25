package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/sensors"
	"TheTinkerDad/sensible/settings"
	"TheTinkerDad/sensible/web/api"
)

var Server *http.Server

func apiHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calling %s...\n", r.URL.Path)

	var result interface{} = nil
	if r.URL.Path == "/api/info" {
		result = api.Info()
	} else if r.URL.Path == "/api/shutdown" {
		result = api.Shutdown(Server)
	} else if r.URL.Path == "/api/pause-mqtt" {
		result = api.PauseMqtt()
	} else if r.URL.Path == "/api/resume-mqtt" {
		result = api.ResumeMqtt()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8090"}
	http.HandleFunc("/api/", apiHandler)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Println("Listening on port 8090...")
	return srv
}

func main() {

	f, err := os.OpenFile("/var/log/sensible/sensible.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Bootstrapping Sensible...")

	settings.EnsureOk()
	mqtt.EnsureOk()
	sensors.EnsureOk()

	serverWaitGroup := &sync.WaitGroup{}
	serverWaitGroup.Add(1)
	Server = startHTTPServer(serverWaitGroup)
	serverWaitGroup.Wait()
}
