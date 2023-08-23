package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/releaseinfo"
	"TheTinkerDad/sensible/sensors"
	"TheTinkerDad/sensible/settings"
	"TheTinkerDad/sensible/web/api"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func setLogLevel(level string) {

	switch level {
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func bootstrap() {

	log.Infof("Bootstrapping Sensible v%s (%s, Commit: %s)", releaseinfo.Version, releaseinfo.BuildTime, releaseinfo.LastCommit)
	settings.Initialize(false)
	settings.GenerateDefaultIfNotExists()
	settings.Load()

	setLogLevel(settings.All.General.LogLevel)

	f, err := os.OpenFile(settings.All.General.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Info("Error opening log file - logging will continue on standard output!")
		log.Infof("Error details: %v", err)
	} else {
		log.SetOutput(f)
	}

	settings.ValidatePluginSettings()

	mqtt.Initialize()
}

func execute() {

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
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		FullTimestamp: true,
	})

	var pversion, phelp, pexample, preset, unregister bool

	flag.BoolVar(&pversion, "v", false, "Show version info.")
	flag.BoolVar(&phelp, "h", false, "Show command line options.")
	flag.BoolVar(&pexample, "g", false, "Generates a sample config file in the working directory.")
	flag.BoolVar(&preset, "r", false, "Reset settings or initialize a fresh install.")
	flag.BoolVar(&unregister, "u", false, "Unregister all sensors from Home Assistant via MQTT.")
	flag.Parse()

	if phelp {
		flag.PrintDefaults()
	} else if pversion {
		fmt.Printf("Sensible v%s (%s, Commit: %s)\n", releaseinfo.Version, releaseinfo.BuildTime, releaseinfo.LastCommit)
	} else if preset {
		log.Info("Setting up defaults...")
		settings.Initialize(false)
		settings.CreateFolders()
		settings.BackupSettingsFile()
		settings.GenerateDefaults()
	} else if pexample {
		log.Info("Generating example configuration...")
		settings.Initialize(true)
		settings.GenerateDefaults()
	} else if unregister {
		bootstrap()
		log.Info("Unregistering all sensors...")
		sensors.UnregisterAllSensors()
	} else {
		bootstrap()
		execute()
	}
}
