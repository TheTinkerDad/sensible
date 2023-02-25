package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"encoding/json"
	"fmt"
	"log"
)

var Paused bool

func RegisterSensor(device DeviceRegistration) {

	payload, err := json.Marshal(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registering new sensor with ID %s: %s\n", device.UniqueId, payload)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, device.UniqueId)
	log.Printf("Configuration topic: %s\n", topic)
	if token := MqttClient.Publish(topic, 1, true, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func RemoveSensor(id string) {

	log.Printf("Unregistering sensor with ID %s", id)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, settings.All.Discovery.DeviceName+"_"+id)
	log.Printf("Configuration topic: %s\n", topic)
	if token := MqttClient.Publish(topic, 1, true, ""); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func SendSensorValue(id string, value string) {

	log.Printf("Sending sensor value for sensor with ID %s: %s", settings.All.Discovery.DeviceName+"_"+id, value)
	topic := fmt.Sprintf("%s/sensor/%s/%s/state", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, settings.All.Discovery.DeviceName+"_"+id)
	log.Printf("State topic: %s\n", topic)
	MqttClient.Publish(topic, 1, false, value)
}

func SendDeviceAvailability(value string) {

	topic := fmt.Sprintf("%s/sensor/%s/availability", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName)
	log.Printf("Sending availability info for device with name %s: %s. Topic: %s", settings.All.Discovery.DeviceName, value, topic)
	MqttClient.Publish(topic, 1, false, value)
}
