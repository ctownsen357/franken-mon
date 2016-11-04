NAME = ctownsend/franken-mon
VERSION = 0.0.1

.PHONY: all clean test install_linux

all:	container linux darwin


container:
	docker build ./golang_build_container/ -t "ctownsend/franken-mon"

test:	container
	docker run --rm -v "$(PWD)/src":/franken-mon ctownsend/franken-mon sh -c \
	"cd /franken-mon && go test"

linux:	container
	docker run --rm -v "$(PWD)/"src:/franken-mon ctownsend/franken-mon sh -c \
	"cd /franken-mon && go build -o franken-mon-linux"

darwin:	container
	docker run --rm -v "$(PWD)/src":/franken-mon --rm -e GOOS=darwin -e GOARCH=amd64 ctownsend/franken-mon sh -c \
	"cd /franken-mon && go build -o franken-mon-darwin"

clean:
	rm -f ./src/franken-mon-linux
	rm -f ./src/franken-mon-darwin
	docker rmi ctownsend/franken-mon

install_linux:
	sudo mkdir /opt/franken-mon
	sudo cp ./src/franken-mon-linux /opt/franken-mon/
	sudo cp ./src/config.json /opt/franken-mon/
	sudo cp franken-mon.service /etc/systemd/system
	sudo systemctl enable franken-mon.service
	sudo systemctl start franken-mon.service
	sudo systemctl status franken-mon.service
