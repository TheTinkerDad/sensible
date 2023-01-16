package sensors

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"log"
	"os/exec"
	"time"
)

// This is temporary stuff (a PoC, really!) and needs to go from here!
// Instead, there needs to be some init code that iterates through
// enabled backend plugins, registers the sensors for each of them.
func registerSensorBootTime() {
	mqtt.RegisterSensor("sensible_boot_time",
		mqtt.DeviceRegistration{
			Name:              "Sensible Boot Time",
			DeviceClass:       "",
			Icon:              "mdi:clock",
			StateTopic:        settings.All.Discovery.Prefix + "/sensor/sensible_boot_time/state",
			UnitOfMeasurement: "",
			ValueTemplate:     "",
			//ValueTemplate:     "{{value_json.value}}",
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
			Name:              "Sensible System Time",
			DeviceClass:       "",
			Icon:              "mdi:clock",
			StateTopic:        settings.All.Discovery.Prefix + "/sensor/sensible_system_time/state",
			UnitOfMeasurement: "",
			ValueTemplate:     "",
			//ValueTemplate:     "{{value_json.value}}",
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

		registerSensorBootTime()
		registerSensorSystemTime()

		for {
			if !mqtt.Paused {
				select {
				case msg := <-SensorUpdater:
					log.Println("Received message", msg)
				default:
					updateSensorBootTime()
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
