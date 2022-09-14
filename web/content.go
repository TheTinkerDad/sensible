package web

import (
	"log"

	rice "github.com/TheTinkerDad/go.rice"
)

// WebContent used to load static web content
var WebContent *rice.Box = nil

func init() {

	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}
	wc, err := conf.FindBox("../data")
	if err != nil {
		wc, err = conf.FindBox("data") // Fallback
		if err != nil {
			log.Fatal(err)
		}
	}
	WebContent = wc

	log.Println("Initializing static web content...")
}

// EnsureOk Checks if the MQTT connection is intact
func EnsureOk() {

}
