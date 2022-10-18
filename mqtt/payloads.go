package mqtt

type DeviceRegistration struct {
	Name              string `json:"name,omitempty"`
	DeviceClass       string `json:"device_class,omitempty"`
	StateTopic        string `json:"state_topic,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string `json:"value_template,omitempty"`
}

type DeviceRemoval struct {
}
