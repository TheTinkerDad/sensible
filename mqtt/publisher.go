package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var Paused bool

func RegisterSensor(device DeviceRegistration) {

	payload, err := json.Marshal(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Registering new sensor with ID %s: %s", device.UniqueId, payload)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, device.UniqueId)
	log.Debugf("Configuration topic: %s", topic)
	if token := MqttClient.Publish(topic, 1, true, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func RemoveSensor(id string) {

	log.Debugf("Unregistering sensor with ID %s", id)
	topic := fmt.Sprintf("%s/sensor/%s/%s/config", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, settings.All.Discovery.DeviceName+"_"+id)
	log.Debugf("Configuration topic: %s", topic)
	if token := MqttClient.Publish(topic, 1, true, ""); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func SendSensorValue(id string, value string) {

	log.Tracef("Sending sensor value for sensor with ID %s: %s", settings.All.Discovery.DeviceName+"_"+id, value)
	topic := fmt.Sprintf("%s/sensor/%s/%s/state", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName, settings.All.Discovery.DeviceName+"_"+id)
	log.Tracef("State topic: %s", topic)
	MqttClient.Publish(topic, 1, false, value)
}

func SendDeviceAvailability(value string) {

	topic := fmt.Sprintf("%s/sensor/%s/availability", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName)
	log.Tracef("Sending availability info for device with name %s: %s. Topic: %s", settings.All.Discovery.DeviceName, value, topic)
	MqttClient.Publish(topic, 1, false, value)
}

func SendAlwaysAvailableMessage() {

	topic := fmt.Sprintf("%s/sensor/%s/always-available", settings.All.Discovery.Prefix, settings.All.Discovery.DeviceName)
	log.Tracef("Sending 'always available' info for device with name %s. Topic: %s", settings.All.Discovery.DeviceName, topic)
	MqttClient.Publish(topic, 1, false, "Online")
}
