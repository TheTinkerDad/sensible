package mqtt

type DeviceMetadata struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
}

type DeviceRegistration struct {
	Name                string         `json:"name,omitempty"`
	DeviceClass         string         `json:"device_class,omitempty"`
	Icon                string         `json:"icon,omitempty"`
	StateTopic          string         `json:"state_topic,omitempty"`
	AvailabilityTopic   string         `json:"availability_topic,omitempty"`
	PayloadAvailable    string         `json:"payload_available,omitempty"`
	PayloadNotAvailable string         `json:"payload_not_available,omitempty"`
	UnitOfMeasurement   string         `json:"unit_of_measurement,omitempty"`
	ValueTemplate       string         `json:"value_template,omitempty"`
	UniqueId            string         `json:"unique_id,omitempty"`
	Device              DeviceMetadata `json:"device,omitempty"`
}

type DeviceRemoval struct {
} 
