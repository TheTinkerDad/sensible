package settings

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type MqttSettings struct {
	Hostname string
	Port     string
	Username string
	Password string
	ClientId string
}

type DiscoverySettings struct {
	DeviceName string
	Prefix     string
}

type AllSettings struct {
	Mqtt      MqttSettings
	Discovery DiscoverySettings
	Plugins   []Plugin
}

type Plugin struct {
	Name     string
	Kind     string
	SensorId string
	Script   string
	Icon     string
}

var All AllSettings

func init() {

	log.Println("Opening configuration file...")
	GenerateDefaultIfNotExists()
	Load()
}

// GenerateDefaultIfNotExists Generates the default configuration file
func GenerateDefaultIfNotExists() {

	if _, err := os.Stat("/etc/sensible/settings.yaml"); errors.Is(err, os.ErrNotExist) {

		log.Println("Config file not found, writing default config...")

		All.Mqtt = MqttSettings{"127.0.0.1", "1883", "", "", "sensible_mqtt_client"}
		All.Discovery = DiscoverySettings{"sensible-1", "homeassistant"}
		All.Plugins = make([]Plugin, 6)
		All.Plugins[0] = Plugin{"Sensible Heartbeat", "internal", "heartbeat", "", "mdi:wrench-check"}
		All.Plugins[1] = Plugin{"Sensible Heartbeat NR", "internal", "heartbeat_NR", "", "mdi:wrench-check"}
		All.Plugins[2] = Plugin{"Sensible Boot Time", "internal", "boot_time", "", "mdi:clock"}
		All.Plugins[3] = Plugin{"Sensible System Time", "internal", "system_time", "", "mdi:clock"}
		All.Plugins[4] = Plugin{"Sensible Root Disk Free", "script", "root_free", "root_free.sh", "mdi:harddisk"}
		All.Plugins[5] = Plugin{"Sensible Host IP Address", "script", "ip_address", "ip_address.sh", "mdi:network"}

		yaml, err := yaml.Marshal(&All)
		if err != nil {
			log.Fatal(err)
		}

		f, err2 := os.Create("/etc/sensible/settings.yaml")
		if err2 != nil {
			log.Fatal(err)
		}
		_, err2 = f.Write(yaml)
		if err2 != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

// Load Loads the current settings
func Load() {

	f, err := os.Open("/etc/sensible/settings.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	raw := make([]byte, fi.Size())
	f.Read(raw)

	err = yaml.Unmarshal(raw, &All)
	if err != nil {
		log.Fatal(err)
	}
}

// EnsureOk Checks if the loaded configuration is intact
func EnsureOk() {

}
