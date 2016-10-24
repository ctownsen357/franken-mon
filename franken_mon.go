package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"
)

/// Configuration is a struct representation of the config.json configuration file.
type Configuration struct {
	Timeout time.Duration
	ActionToMonitor string
	ContainerNamesToMonitor	map[string]bool
	CommandTemplates    []string
}

/// GetConfig returns a Configuration struct from the applicaton config.json file
func GetConfig() Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Println("There was an error trying to parse the configuration file config.json:")
		log.Fatal(err)
	}

	return conf
}

// ParseCommandTemplate takes a command template from the configuration file and an event message
// and parses the template into a command string
func ParseCommandTemplate(cmdTemplate string, msg *docker.APIEvents) (string, error){
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
	conf := GetConfig()

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

	timeout := time.After(conf.Timeout * time.Second)

	for {
		select {
		case msg := <-listener:
			if (msg.Action == conf.ActionToMonitor && len(conf.ContainerNamesToMonitor) == 0) || (len(conf.ContainerNamesToMonitor) > 0 && conf.ContainerNamesToMonitor[msg.Actor.Attributes["name"]]){
				log.Println(msg)
				log.Println(msg.ID, msg.Action, msg.Actor.Attributes["name"])
				conf = GetConfig() //re-obtaining the command list to obtain any changes or additions since last start event

				//create a command template based on the configuration file
				//pass the ID from the message start event to the template
				for _, cmdTemplate := range conf.CommandTemplates {
					parsedCmd, err := ParseCommandTemplate(cmdTemplate, msg)
					if err != nil {
						log.Fatal(err)
					}
					// I'd like to explore the REST API or the go-dockerclient exec options
					// further but I haven't used either and am implementing this as quickly as
					// I can to finish the excercise in a timely manner.
					// /bin/sh should be available in most containers / distros but this could also
					// be a config file option or a per command option if necessary (and one didn't implement this via the RESt or go-dockerclient APIs)
					cmd := exec.Command("/bin/sh", "-c", parsedCmd)

					//pipe the cmd stdout & stderr to the log so we have a record of what happened
					stdout, err := cmd.StdoutPipe()
					if err != nil {
						log.Fatal(err)
					}
					stderr, err := cmd.StderrPipe()
					if err != nil {
						log.Fatal(err)
					}
					multi := io.MultiReader(stdout, stderr)

					err = cmd.Start()

					if err != nil {
						log.Fatal(err)
					}
					in := bufio.NewScanner(multi)
					for in.Scan() {
						log.Printf(in.Text())
					}
				}
			}
		case <-timeout:
			log.Println("Timeout period reached without seeing any Docker events.")
			return
		}
	}

}
