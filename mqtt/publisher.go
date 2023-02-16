package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"encoding/json"
	"fmt"
	"log"
)

var Paused bool

func RegisterSensor(id string, device DeviceRegistration) {

	payload, err := json.Marshal(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registering new sensor with ID %s: %s\n", id, payload)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, id)
	log.Printf("Configuration topic: %s\n", topic)
	if token := MqttClient.Publish(topic, 1, true, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func RemoveSensor(id string) {

	log.Printf("Unregistering sensor with ID %s", id)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, id)
	log.Printf("Configuration topic: %s\n", topic)
	MqttClient.Publish(topic, 1, true, nil)
}

func SendSensorValue(id string, value string) {
	log.Printf("Sending sensor value for sensor with ID %s: %s", id, value)
	log.Printf(" --> %s/sensor/%s/%s/state\n",
		settings.All.Discovery.Prefix,
		settings.All.Discovery.DeviceName,
		id)
	MqttClient.Publish(fmt.Sprintf("%s/sensor/%s/%s/state",
		settings.All.Discovery.Prefix,
		settings.All.Discovery.DeviceName,
		id), 1, false, value)
}

func SendDeviceAvailability(value string) {
	log.Printf("Sending availability info for device with name %s: %s", settings.All.Discovery.DeviceName, value)
	log.Printf(" --> %s/sensor/%s/availability\n",
		settings.All.Discovery.Prefix,
		settings.All.Discovery.DeviceName)
	MqttClient.Publish(fmt.Sprintf("%s/sensor/%s/availability",
		settings.All.Discovery.Prefix,
		settings.All.Discovery.DeviceName), 1, false, value)
}
