package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/releaseinfo"
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

func execute() {
	log.Printf("Bootstrapping Sensible v%s (%s)\n", releaseinfo.Version, releaseinfo.BuildTime)
	settings.EnsureOk()

	f, err := os.OpenFile(settings.All.General.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening log file - logging will continue on standard output!")
		log.Printf("Error details: %v\n", err)
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	mqtt.EnsureOk()
	sensors.EnsureOk()

	serverWaitGroup := &sync.WaitGroup{}
	serverWaitGroup.Add(1)
	Server = startHTTPServer(serverWaitGroup)
	serverWaitGroup.Wait()
}

func main() {

	log.SetOutput(os.Stdout)

	var pversion bool
	var phelp bool
	var preset bool

	flag.BoolVar(&pversion, "v", false, "Show version info.")
	flag.BoolVar(&phelp, "h", false, "Show command line options.")
	flag.BoolVar(&preset, "r", false, "Reset settings or initialize a fresh install.")
	flag.Parse()

	if phelp {
		flag.PrintDefaults()
	} else if pversion {
		fmt.Printf("Sensible v%s (%s)\n", releaseinfo.Version, releaseinfo.BuildTime)
	} else if preset {
		log.Println("Setting up defaults...")
		settings.CreateFolders()
		settings.BackupSettingsFile()
		settings.GenerateDefaults()
	} else {
		execute()
	}
}
