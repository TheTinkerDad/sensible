package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"encoding/json"
	"fmt"
	"log"
)

func RegisterSensor(id string, device DeviceRegistration) {

	payload, err := json.Marshal(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registering new device with ID %s: %s\n", id, payload)
	topic := fmt.Sprintf("%s/sensor/%s/config", settings.All.Discovery.Prefix, id)
	log.Printf("Configuration topic: %s\n", topic)
	if token := MqttClient.Publish(topic, 1, true, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func RemoveSensor(id string) {

	log.Printf("Unregistering device with ID %s\n", id)
	MqttClient.Publish(fmt.Sprintf("%s/sensor/%s/config", settings.All.Discovery.Prefix, id), 1, true, DeviceRemoval{})
}

func SendSensorValue(id string, value string) {
	log.Printf("Sending sensor value for device with ID %s: %s\n", id, value)
	MqttClient.Publish(fmt.Sprintf("%s/sensor/%s/state", settings.All.Discovery.Prefix, id), 1, false, value)
}
