# franken-mon
![franken-mon](images/alive.jpg "It's alive!")

franken-mon was an October excercise, hence the Halloween inspired name, to create a service in Go that could be used to monitor Docker restart events and execute commands against the restarted container.  For example, a container that needs to have any stateful type of information set on restart; something like this could be used to accomplish that objective. This was to be done using a container for the build steps so that one only needed to have Docker and GNU Make installed - not Go. At the end of the day this was a fun and worthwhile excercise.  I'll likely develop this further as I have time.

Alternatively, one could implement a simple shell script and accomplish nearly the same thing with the following:

```{bash}
echo ' pwd' > cmd.txt && docker events | awk '/container restart/{system("echo docker exec " $4 " $(cat cmd.txt) | bash -")}'
```

That command pipes the pwd command into a file and then pipes the stream from docker events into awk which is searching for the restart event. When the restart event is encountered it executes the arbitrary command from the text file against the restarted container. The command in the text file could be replaced with any desired command.


Where is the fun in that; let's write some code, make it testable, and explores doing the same thing running a Go binary as a systemd service.

### The request:
Write a configurable service in Go that can be started by running a compiled binary. The service should monitor the Docker API for restart events and run an arbitrary command in response to that event. The arbitrary command should be supplied via the config file and should allow a template like {{ .ID }} so the user can run commands against the restarted container.

Deliverables: A golang tool that reads an arbitrary command like docker exec {{ .ID }} pwd from a config file and runs it against containers that have recently been restarted. The code should endeavor to only run the command once per restart. A Dockerfile for a container that when run with a mounted directory will compile the Golang code and make the binary available to the host. This should compile versions for both Darwin and Linux. When running, the monitor should log(to a file) the container ID, Name of any qualifying events along with the outcome of the command it runs in response to the event.

Bonus: A makefile with helpers for compiling the golang code in a container. The ability to provide the names of containers the user wants monitored The ability to run the watching service in the background on OSX Tests


### Installing / building: ###
``` {bash}
git clone https://github.com/ctownsen357/franken-mon.git
cd franken-mon
make # This builds binaries for both Linux and OS X; you could optionally specify the platform like: make linux or make darwin
make install_linux
```


### Monitoring with journalctl: ###
```
sudo journalctl -u franken-mon -f
```

### Testing: ###
Once you have the service installed or running via the command line.  Open up another terminal and restart a docker container (**with the same name as the container in your config.json file if you choose to limit by container name **).
