package sensors

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"bytes"
	"log"
	"os/exec"
	"time"

	pipe "github.com/TheTinkerDad/go.pipe"
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

	dr := mqtt.DeviceRegistration{
		Name:                name,
		DeviceClass:         "",
		Icon:                icon,
		StateTopic:          settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/" + settings.All.Discovery.DeviceName + "_" + id + "/state",
		AvailabilityTopic:   settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability",
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

func updateSensorHeartbeat() {

	mqtt.SendSensorValue("heartbeat", "ONLINE")
}

func updateSensorSystemTime() {

	now := time.Now()
	mqtt.SendSensorValue("system_time", string(now.Format("2006-01-02 15:04:05")))
}

func updateSensorBootTime() {

	// TODO: Find an OS-agnostic solution for this!
	out, err := exec.Command("uptime", "-s").Output()
	if err != nil {
		log.Println(err)
		out = []byte("Unavailable")
	}
	mqtt.SendSensorValue("boot_time", string(out))
}

// This updates sensors based on scripts

func updateSensorWithScript(p settings.Plugin) {

	log.Printf("Executing %s%s\n", settings.All.General.ScriptLocation, p.Script)
	// Using pipe here looks like an overkill, but can be useful later...
	var b bytes.Buffer
	if err := pipe.Command(&b,
		exec.Command("sh", "-c", settings.All.General.ScriptLocation+p.Script),
	); err != nil {
		log.Fatal(err)
	}
	mqtt.SendSensorValue(p.SensorId, b.String())
}

var SensorUpdater chan string

// EnsureOk Triggers init() which starts registering the devices
func EnsureOk() {

	go func() {

		for _, p := range settings.All.Plugins {
			mqtt.RegisterSensor(getSensorMetaData(p.SensorId, p.Name, p.Icon, p.UnitOfMeasurement))
		}

		for {
			if !mqtt.Paused {
				select {
				case msg := <-SensorUpdater:
					log.Println("Received message", msg)
				default:
					mqtt.SendDeviceAvailability("Online")
					for _, p := range settings.All.Plugins {
						switch p.Kind {
						case "internal":
							//TODO: This should be reflection based!
							switch p.SensorId {
							case "heartbeat":
								updateSensorHeartbeat()
							case "heartbeat_NR":
								// This is a sensor we reserved for NodeRed, so we don't update it here.
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
