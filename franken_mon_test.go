package main

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestGetConfig(t *testing.T) {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	expected := "restart"
	expect(t, config.ActionToMonitor, expected)
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
		log.Fatal(err)
	}
	expected := "docker exec 5745704abe9caa5 pwd"
	expect(t, rslt, expected)
}

func TestExecuteCommand(t *testing.T) {
	rslt, err := ExecuteCommand("echo 'testing 1,2,3...'")
	if err != nil {
		log.Fatal("Error trying to execute command:", err)
	}

	expected := "testing 1,2,3..."

	if !strings.Contains(rslt, expected) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), rslt, reflect.TypeOf(rslt))
	}

}
