package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var MqttClient mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debugf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Warnf("Connect lost: %v", err)
}

// Initialize Checks if the MQTT connection is intact
func Initialize() {

	log.Infof("Connecting to MQTT broker at %s:%s...", settings.All.Mqtt.Hostname, settings.All.Mqtt.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", settings.All.Mqtt.Hostname, settings.All.Mqtt.Port))
	opts.SetClientID(settings.All.Mqtt.ClientId)
	opts.SetUsername(settings.All.Mqtt.Username)
	opts.SetPassword(settings.All.Mqtt.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetWill(settings.All.Discovery.Prefix+"/sensor/"+settings.All.Discovery.DeviceName+"/availability", "Offline", 1, false)
	MqttClient = mqtt.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
