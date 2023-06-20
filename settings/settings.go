package settings

import (
	"TheTinkerDad/sensible/utility"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GeneralSettings struct {
	Logfile        string
	LogLevel       string
	ScriptLocation string
}

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

type ApiSettings struct {
	Enabled bool
	Port    int
	Token   string
}

type AllSettings struct {
	General   GeneralSettings
	Mqtt      MqttSettings
	Discovery DiscoverySettings
	Api       ApiSettings
	Plugins   []Plugin
}

type Plugin struct {
	Name              string
	Kind              string
	SensorId          string
	Script            string
	UnitOfMeasurement string
	Icon              string
	UpdateInterval    int64
	Value             string
}

var PluginDefaults = Plugin{"!", "", "!", "", "", "mdi:wrench-check", 10, ""}

var All AllSettings

var settingsFolder string = "/etc/sensible/"
var settingsFile string = "/etc/sensible/settings.yaml"

// Backs up the existing settings file - if there's any
func BackupSettingsFile() {

	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {
		return
	} else {
		utility.Copy(settingsFile, settingsFile+".bkp")
	}
}

// Generates the default configuration file
func GenerateDefaults() {

	// We're running Sensible as a non-root user
	if os.Getuid() == 0 {
		All.General = GeneralSettings{"/var/log/sensible/sensible.log", "info", "/etc/sensible/scripts/"}
	} else {
		All.General = GeneralSettings{settingsFolder + "log/sensible.log", "info", settingsFolder + "scripts/"}
	}
	All.Mqtt = MqttSettings{"127.0.0.1", "1883", "", "", "sensible_mqtt_client"}
	All.Discovery = DiscoverySettings{"sensible-demo", "homeassistant"}
	All.Api = ApiSettings{Port: 8090, Enabled: false, Token: utility.NewRandomUUID()}
	All.Plugins = make([]Plugin, 8)
	All.Plugins[0] = Plugin{"Heartbeat", "internal", "heartbeat", "", "", "mdi:wrench-check", 10, ""}
	All.Plugins[1] = Plugin{"Boot Time", "internal", "boot_time", "", "", "mdi:clock", 120, ""}
	All.Plugins[2] = Plugin{"System Time", "internal", "system_time", "", "", "mdi:clock", 10, ""}
	All.Plugins[3] = Plugin{"Version", "internal", "version", "", "", "mdi:wrench-check", 3600, ""}
	All.Plugins[4] = Plugin{"Root Disk Free", "script", "root_free", "root_free.sh", "GB", "mdi:harddisk", 60, ""}
	All.Plugins[5] = Plugin{"Host IP Address", "script", "ip_address", "ip_address.sh", "", "mdi:network", 600, ""}
	All.Plugins[6] = Plugin{"Hostname", "script", "hostname", "hostname.sh", "", "mdi:network", 1200, ""}
	All.Plugins[7] = Plugin{"Platform", "script", "platform", "platform.sh", "", "mdi:wrench-check", 1200, ""}

	yaml, err := yaml.Marshal(&All)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(settingsFile, yaml, 0600)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Infof("Successfully written default config file as %s!", settingsFile)

	/*
	   f, err2 := os.OpenFile(settingsFile, os.O_CREATE, 0600)

	   	if err2 != nil {
	   		log.Fatal(err)
	   	}

	   _, err2 = f.Write(yaml)

	   	if err2 != nil {
	   		log.Fatal(err)
	   	}

	   f.Close()
	*/
}

// CreateFolders Creates the default folders used by Sensible
func CreateFolders() {

	log.Info("Creating default folders...")

	// We're running Sensible as a non-root user
	if os.Getuid() == 0 {
		utility.CreateFolder("/etc/sensible/scripts/")
		utility.CreateFolder("/var/log/sensible")
	} else {
		utility.CreateFolder(settingsFolder + "scripts")
		utility.CreateFolder(settingsFolder + "log")
	}

}

// GenerateDefaultIfNotExists Generates the default configuration file
func GenerateDefaultIfNotExists() {

	log.Debug("Opening configuration file...")
	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {

		log.Warn("Config file not found, writing default config...")
		GenerateDefaults()
	}
}

// Load Loads the current settings
func Load() {

	f, err := os.Open(settingsFile)
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

func hasValue(settingType string, settingVal interface{}) bool {

	return (settingType == "string" && settingVal.(string) != "") ||
		(settingType == "int64" && settingVal.(int64) != 0)
}

func failWithMissingValue(pluginIndex int, fieldName string) {
	log.Errorf("Configuration field %s of plugin #%d must have a value!", fieldName, pluginIndex)
	os.Exit(-1)
}

func failWithInvalidValue(pluginIndex int, fieldName string, fieldValue string) {
	log.Errorf("Configuration field %s of plugin #%d has an invalid value: %s", fieldName, pluginIndex, fieldValue)
	os.Exit(-1)
}

func ValidatePluginSettings() {

	d := reflect.ValueOf(PluginDefaults)
	pluginType := d.Type()
	pluginIndex := 1
	for _, p := range All.Plugins {
		v := reflect.ValueOf(p)
		for i := 0; i < v.NumField(); i++ {
			fieldName := pluginType.Field(i).Name
			fieldType := pluginType.Field(i).Type.Name()
			fieldValue := v.Field(i).Interface()
			fieldDefault := d.Field(i).Interface()
			if !hasValue(fieldType, fieldValue) {
				if fieldDefault == "!" {
					failWithMissingValue(pluginIndex, fieldName)
				} else if hasValue(fieldType, fieldDefault) {
					log.Debugf("Plugin %d. missing config variable %s of type %s!", pluginIndex, fieldName, fieldType)
					if fieldType == "string" {
						log.Debugf("Using default value '%s' for configuration field %s of plugin #%d.", fieldDefault.(string), fieldName, pluginIndex)
						reflect.ValueOf(&p).Elem().Field(i).SetString(fieldDefault.(string))
					} else if fieldType == "int64" {
						log.Debugf("Using default value '%d' for configuration field %s of plugin #%d.", fieldDefault.(int64), fieldName, pluginIndex)
						reflect.ValueOf(&p).Elem().Field(i).SetInt(fieldDefault.(int64))
					}
				}
			} else {
				validateFieldValue(pluginIndex, fieldName, fieldType, fieldValue)
			}
		}
		pluginIndex++
	}
}

func validateFieldValue(pluginIndex int, fieldName string, fieldType string, fieldValue interface{}) {

	switch fieldName {
	case "Kind":
		switch fieldValue {
		case "internal":
		case "script":
		case "fixed":
			return
		default:
			failWithInvalidValue(pluginIndex, fieldName, fmt.Sprintf("%v", fieldValue))
		}
	}
}

// Initializes The location of the settings based on the detected user
func Initialize() {

	// We're running Sensible as a non-root user
	if os.Getuid() != 0 {
		settingsFolder = os.Getenv("HOME") + "/.sensible/"
		settingsFile = settingsFolder + "settings.yaml"
	}
}
