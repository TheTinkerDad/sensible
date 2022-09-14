package api

import (
	"log"
	"os/exec"
)

// AnsibleInfoStruct Holds data about the local Ansible installation
type AnsibleInfoStruct struct {
	Version string
}

// AnsibleInfo returns a JSON that provides Ansible info
func AnsibleInfo() AnsibleInfoStruct {

	out, err := exec.Command("ansible", "--version").Output()
	if err != nil {
		log.Println(err)
		out = []byte("Ansible not installed.")
	}

	return AnsibleInfoStruct{Version: string(out)}
}
