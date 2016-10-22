package main

import (
	"testing"
	"github.com/fsouza/go-dockerclient"
)

func TestGetConfig(t *testing.T) {
	var config Configuration = GetConfig() //explicitly declaring the expected type as part of the test
	if config.Timeout <= 0 {
		t.Fatal("Timout must be greater than 0.", config.Timeout)
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

				Time:     1442421716,
				TimeNano: 1442421716983607193,
			},
		}

	rslt, err := ParseCommandTemplate(cmdTemplate, events[0])
	if err != nil {
		t.Fatal(err)
	}
	expected := "docker exec 5745704abe9caa5 pwd"
	if rslt != expected {
		t.Fatal("The command did not parse correctly.  Expected %v but got %v", expected, rslt)
	}
}
