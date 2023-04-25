package sensors

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"bytes"
	"os/exec"
	"strings"
	"sync"
	"time"

	pipe "github.com/TheTinkerDad/go.pipe"
	log "github.com/sirupsen/logrus"
)

func getDeviceMetaData() mqtt.DeviceMetadata {

	dmd := mqtt.DeviceMetadata{
		Name:         settings.All.Discovery.DeviceName,
		Manufacturer: "TheTinkerDad",
		Model:        "Sensible-Sensor",
	}
	dmd.Identifiers = make([]string, 1)
	dmd.Identifiers[0] = ("sensible-" + settings.All.Discovery.DeviceName)
	return dmd
}

func getSensorMetaData(id string, name string, icon string, unit string) mqtt.DeviceRegistration {

	var deviceClass, stateTopic, availabilityTopic string
	if id == "heartbeat" {
		stateTopic = settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability"
		availabilityTopic = settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/always-available"
	} else {
		stateTopic = settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/" + settings.All.Discovery.DeviceName + "_" + id + "/state"
		availabilityTopic = settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability"
	}

	dr := mqtt.DeviceRegistration{
		Name:                name,
		DeviceClass:         deviceClass,
		Icon:                icon,
		StateTopic:          stateTopic,
		AvailabilityTopic:   availabilityTopic,
		PayloadAvailable:    "Online",
		PayloadNotAvailable: "Offline",
		UnitOfMeasurement:   unit,
		ValueTemplate:       "",
		//ValueTemplate:     "{{value_json.value}}",
		UniqueId: settings.All.Discovery.DeviceName + "_" + id,
		Device:   getDeviceMetaData(),
	}
	return dr
}

// These below methods are for updating the simple internal sensors we currently have

func updateSensorSystemTime() {

	now := time.Now()
	mqtt.SendSensorValue("system_time", string(now.Format("2006-01-02 15:04:05")))
}

func updateSensorBootTime() {

	// TODO: Find an OS-agnostic solution for this!
	out, err := exec.Command("uptime", "-s").Output()
	if err != nil {
		log.Warn(err)
		out = []byte("Unavailable")
	}
	value := strings.TrimSuffix(string(out), "\n")
	mqtt.SendSensorValue("boot_time", value)
}

// This updates sensors based on scripts

func updateSensorWithScript(p settings.Plugin) {

	log.Tracef("Executing %s%s", settings.All.General.ScriptLocation, p.Script)
	// Using pipe here looks like an overkill, but can be useful later...
	var b bytes.Buffer
	if err := pipe.Command(&b,
		exec.Command("sh", "-c", settings.All.General.ScriptLocation+p.Script),
	); err != nil {
		log.Warn(err)
	}
	value := strings.TrimSuffix(b.String(), "\n")
	mqtt.SendSensorValue(p.SensorId, value)
}

var SensorUpdater chan string

// UnregisterAllSensors Sends MQTT messages to HA to deregister all sensors
func UnregisterAllSensors() {

	for _, p := range settings.All.Plugins {
		mqtt.RemoveSensor(p.SensorId)
	}
}

// StartProcessing Starts the loop to process and send sensor data
func StartProcessing(wg *sync.WaitGroup) {

	go func() {

		defer wg.Done()

		for _, p := range settings.All.Plugins {
			mqtt.RegisterSensor(getSensorMetaData(p.SensorId, p.Name, p.Icon, p.UnitOfMeasurement))
		}

		log.Info("Entering MQTT message processing loop...")

		for {
			if !mqtt.Paused {
				select {
				case msg := <-SensorUpdater:
					log.Tracef("Received message %s", msg)
				default:
					mqtt.SendAlwaysAvailableMessage()
					mqtt.SendDeviceAvailability("Online")
					for _, p := range settings.All.Plugins {
						switch p.Kind {
						case "internal":
							//TODO: This should be reflection based!
							switch p.SensorId {
							case "boot_time":
								updateSensorBootTime()
							case "system_time":
								updateSensorSystemTime()
							default:
							}
						case "script":
							updateSensorWithScript(p)
						default:
						}
					}
				}
			}
			// TODO: This sould be removed and update periodicity should be configurable on a per-sensor basis
			time.Sleep(10 * time.Second)
		}
	}()
}
