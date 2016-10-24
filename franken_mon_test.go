package main

import (
	"strings"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func TestGetConfig(t *testing.T) {
	config, err := GetConfig()
	if err != nil {
		t.Fatal("Error was not nil:", err)
	}

	expected := "restart"
	if config.ActionToMonitor != expected {
		t.Errorf("Configuration file is not returning expected value.  Expected %s but received %s ", expected, config.ActionToMonitor)
	}
}

func TestParseCommandTemplate(t *testing.T) {
	cmdTemplate := "docker exec {{.ID}} pwd"
	events := []*docker.APIEvents{
		{
			Action: "start",
			Type:   "container",
			Actor: docker.APIActor{
				ID: "5745704abe9caa5",
				Attributes: map[string]string{
					"image": "alpine",
				},
			},

			Status: "start",
			ID:     "5745704abe9caa5",
			From:   "alpine",
		},
	}

	rslt, err := ParseCommandTemplate(cmdTemplate, events[0])
	if err != nil {
		t.Fatal(err)
	}
	expected := "docker exec 5745704abe9caa5 pwd"
	if rslt != expected {
		t.Errorf("The command did not parse correctly.  Expected %v but got %v", expected, rslt)
	}
}

func TestExecuteCommand(t *testing.T) {
	rslt, err := ExecuteCommand("echo 'testing 1,2,3...'")
	if err != nil {
		t.Fatal("Error trying to execute command:", err)
	}

	expected := "testing 1,2,3..."

	if !strings.Contains(rslt, expected) {
		t.Errorf("The command did not run correctly.  Expected %v but got %v", expected, rslt)
	}

}
