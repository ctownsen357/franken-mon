package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/fsouza/go-dockerclient"
)

// Configuration is a struct representation of the config.json configuration file.
type Configuration struct {
	ActionToMonitor         string
	ContainerNamesToMonitor map[string]bool
	CommandTemplates        []string
}

// GetConfig returns a Configuration struct from the applicaton config.json file
func GetConfig() (Configuration, error) {
	conf := Configuration{}
	file, err := os.Open("config.json")
	if err != nil {
		return conf, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

// ExecuteCommand takes a complete command as a string and executes it in the shell environment; intended to execute Docker commands
// but would execute any valid command.
func ExecuteCommand(parsedCmd string) (string, error) {
	// /bin/sh should be available in most containers / distros but this could also
	// be a config file option or a per command option if necessary (and one didn't implement this via the REST or go-dockerclient APIs)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", parsedCmd)

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Start()
	if err != nil {
		return "", err
	}
	cmd.Wait()
	return outbuf.String() + errbuf.String(), err
}

// ParseCommandTemplate takes a command template from the configuration file and an event message
// and parses the template into a command string
func ParseCommandTemplate(cmdTemplate string, msg *docker.APIEvents) (string, error) {
	tmpl, err := template.New("cmd").Parse(cmdTemplate)
	if err != nil {
		return "", err
	}
	var cmdBytes bytes.Buffer
	err = tmpl.Execute(&cmdBytes, msg)
	if err != nil {
		return "", err
	}

	return cmdBytes.String(), nil
}

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	listener := make(chan *docker.APIEvents)
	err = client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = client.RemoveEventListener(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	conf, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	if len(conf.ContainerNamesToMonitor) > 0{
		log.Printf("Waiting for %s events for the following container(s): %v; if you want to monitor all containers please delete the container names from the ContainerNamesToMonitor collection in config.json.  You do NOT need to restart the service for the change to take effect.", conf.ActionToMonitor, conf.ContainerNamesToMonitor)
	}
	for {
		select {
		case msg := <-listener:
			if msg.Action == conf.ActionToMonitor {
				conf, err := GetConfig() //re-loading the config/command list to obtain any changes or additions since last start event
				if err != nil {
					log.Fatal(err)
				}
				if len(conf.ContainerNamesToMonitor) == 0 || conf.ContainerNamesToMonitor[msg.Actor.Attributes["name"]] {
					log.Println(msg.ID, msg.Action, msg.Actor.Attributes["name"])

					//create a command template based on the configuration file
					//pass the ID from the message start event to the template
					for _, cmdTemplate := range conf.CommandTemplates {
						parsedCmd, err := ParseCommandTemplate(cmdTemplate, msg)
						if err != nil {
							log.Fatal(err)
						}
						rslt, err := ExecuteCommand(parsedCmd)
						if err != nil {
							log.Fatal(err)
						}
						log.Println(rslt)
					}
				}
			}
		}
	}

}
