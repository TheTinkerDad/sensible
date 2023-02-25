package mqtt

import (
	"TheTinkerDad/sensible/settings"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttClient mqtt.Client

func init() {

	log.Printf("Connecting to MQTT broker at %s:%s...\n", settings.All.Mqtt.Hostname, settings.All.Mqtt.Port)
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
	MqttClient = mqtt.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v\n", err)
}

// EnsureOk Checks if the MQTT connection is intact
func EnsureOk() {

}
