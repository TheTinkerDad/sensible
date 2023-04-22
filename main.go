package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/releaseinfo"
	"TheTinkerDad/sensible/sensors"
	"TheTinkerDad/sensible/settings"
	"TheTinkerDad/sensible/web/api"
)

func execute() {

	log.Printf("Bootstrapping Sensible v%s (%s)\n", releaseinfo.Version, releaseinfo.BuildTime)
	settings.Initialize()

	f, err := os.OpenFile(settings.All.General.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening log file - logging will continue on standard output!")
		log.Printf("Error details: %v\n", err)
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	mqtt.Initialize()

	funcWaitGroup := &sync.WaitGroup{}

	if settings.All.Api.Enabled {
		funcWaitGroup.Add(1)
		api.StartApiServer(funcWaitGroup)
	}

	funcWaitGroup.Add(1)
	sensors.StartProcessing(funcWaitGroup)

	funcWaitGroup.Wait()
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
