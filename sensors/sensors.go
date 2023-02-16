package sensors

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"log"
	"os/exec"
	"time"
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

// This is temporary stuff (a PoC, really!) and needs to go from here!
// Instead, there needs to be some init code that iterates through
// enabled backend plugins, registers the sensors for each of them.

func registerSensorHeartbeat() {
	mqtt.RegisterSensor("sensible_heartbeat",
		mqtt.DeviceRegistration{
			Name:                "Sensible Heartbeat",
			DeviceClass:         "",
			Icon:                "mdi:wrench-check",
			StateTopic:          settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/sensible_heartbeat/state",
			AvailabilityTopic:   settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability",
			PayloadAvailable:    "Online",
			PayloadNotAvailable: "Offline",
			UnitOfMeasurement:   "",
			ValueTemplate:       "",
			//ValueTemplate:     "{{value_json.value}}",
			UniqueId: settings.All.Discovery.DeviceName + "_heartbeat",
			Device:   getDeviceMetaData(),
		})
}
func registerSensorHeartbeatNR() {
	mqtt.RegisterSensor("sensible_heartbeat_NR",
		mqtt.DeviceRegistration{
			Name:                "Sensible Heartbeat_NR",
			DeviceClass:         "",
			Icon:                "mdi:wrench-check",
			StateTopic:          settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/sensible_heartbeat_NR/state",
			AvailabilityTopic:   settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability",
			PayloadAvailable:    "Online",
			PayloadNotAvailable: "Offline",
			UnitOfMeasurement:   "",
			ValueTemplate:       "",
			//ValueTemplate:     "{{value_json.value}}",
			UniqueId: settings.All.Discovery.DeviceName + "_heartbeat_nr",
			Device:   getDeviceMetaData(),
		})
}

func updateSensorHeartbeat() {

	mqtt.SendSensorValue("sensible_heartbeat", "ONLINE")
}

func registerSensorBootTime() {
	mqtt.RegisterSensor("sensible_boot_time",
		mqtt.DeviceRegistration{
			Name:                "Sensible Boot Time",
			DeviceClass:         "",
			Icon:                "mdi:clock",
			StateTopic:          settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/sensible_boot_time/state",
			AvailabilityTopic:   settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability",
			PayloadAvailable:    "Online",
			PayloadNotAvailable: "Offline",
			UnitOfMeasurement:   "",
			ValueTemplate:       "",
			//ValueTemplate:     "{{value_json.value}}",
			UniqueId: settings.All.Discovery.DeviceName + "_boottime",
			Device:   getDeviceMetaData(),
		})
}

func updateSensorBootTime() {

	out, err := exec.Command("uptime", "-s").Output()
	if err != nil {
		log.Println(err)
		out = []byte("Unavailable")
	}
	mqtt.SendSensorValue("sensible_boot_time", string(out))
}

func registerSensorSystemTime() {
	mqtt.RegisterSensor("sensible_system_time",
		mqtt.DeviceRegistration{
			Name:                "Sensible System Time",
			DeviceClass:         "",
			Icon:                "mdi:clock",
			StateTopic:          settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/sensible_system_time/state",
			AvailabilityTopic:   settings.All.Discovery.Prefix + "/sensor/" + settings.All.Discovery.DeviceName + "/availability",
			PayloadAvailable:    "Online",
			PayloadNotAvailable: "Offline",
			UnitOfMeasurement:   "",
			ValueTemplate:       "",
			//ValueTemplate:     "{{value_json.value}}",
			UniqueId: settings.All.Discovery.DeviceName + "_systemtime",
			Device:   getDeviceMetaData(),
		})
}

func updateSensorSystemTime() {

	out, err := exec.Command("date", "+%Y-%m-%d %T").Output()
	if err != nil {
		log.Println(err)
		out = []byte("Unavailable")
	}
	mqtt.SendSensorValue("sensible_system_time", string(out))
}

var SensorUpdater chan string

func init() {

	go func() {

		//		mqtt.RemoveSensor("sensible_os_uptime")
		//		mqtt.RemoveSensor("sensible_heartbeat")
		//		mqtt.RemoveSensor("sensible_heartbeat_NR")
		//		mqtt.RemoveSensor("sensible_boot_time")
		//		mqtt.RemoveSensor("sensible_system_time")

		//		registerSensorHeartbeat()
		//		registerSensorHeartbeatNR()
		//		registerSensorBootTime()
		registerSensorSystemTime()

		for {
			if !mqtt.Paused {
				select {
				case msg := <-SensorUpdater:
					log.Println("Received message", msg)
				default:
					mqtt.SendDeviceAvailability("Online")
					//					updateSensorHeartbeat()
					//					updateSensorBootTime()
					updateSensorSystemTime()
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

// EnsureOk Triggers init() which starts registering the devices
func EnsureOk() {

}
