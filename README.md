# franken-mon

### Assumptions: ###
Ordinarily I'd ask for more information to define the requirements a little more completely in the areas where I've made assumptions.   Since the intent was primarily to demonstrate a level of proficiency with Docker and Go I got started with the following assumptions:
- Creating the Linux service as a systemd service was OK since we discussed CoreOS, Kubernetes, etc
- One Docker event only (I thought it would be interesting to make it monitor any number of events and container name combination but I wanted to get to a stopping point.)
- I included an example of how to compile for OS X but I figured it was OK to omit the Mac OS X service creation.  I do have access to a Mac Book Pro but most of my machines at home and at work are Linux machines so I haven't had the need to create services for OS X.  I'm confident I could easily do that but again this is a question I'd get answered in further detail to make sure it was done as needed / required.  
- I chose to have the Makefile install the systemd service to /opt/franken-mon .  In a production environment I'd use something like Ansible to deploy and make the install directory and service user account configurable.  


### Installing / building: ###
``` {bash}
git clone https://github.com/drud/ctownsend.git
cd drud/ctownsend/franken-mon
make # This builds binaries for both Linux and OS X; you could optionally specify the platform like: make linux or make darwin
make install_linux
```


### Monitoring with journalctl: ###
```
sudo journalctl -u franken-mon -f
```

### Testing: ###
Once you have the service installed or running via the command line.  Open up another terminal and restart a docker container (**with the same name as the container in your config.json file if you choose to limit by container name **).
