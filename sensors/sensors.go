package sensors

import (
	"TheTinkerDad/sensible/mqtt"
	"TheTinkerDad/sensible/settings"
	"log"
	"os/exec"
)

func init() {

	// This is temporary stuff (a PoC, really!) and needs to go from here!
	// Instead, there needs to be some init code that iterates through
	// enabled backend plugins, registers the sensors for each of them.
	mqtt.RegisterSensor("sensible_os_uptime",
		mqtt.DeviceRegistration{
			Name:              "OS Uptime",
			DeviceClass:       "",
			StateTopic:        settings.All.Discovery.Prefix + "/sensor/sensible_os_uptime/state",
			UnitOfMeasurement: "",
			ValueTemplate:     "",
			//ValueTemplate:     "{{value_json.value}}",
		})

	out, err := exec.Command("uptime", "-s").Output()
	if err != nil {
		log.Println(err)
		out = []byte("Unavailable")
	}
	mqtt.SendSensorValue("sensible_os_uptime", string(out))
}

// EnsureOk Triggers init() which starts registering the devices
func EnsureOk() {

}
