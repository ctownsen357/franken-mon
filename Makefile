NAME = ctownsend/franken-mon
VERSION = 0.0.1

.PHONY: all build clean

all:	build

build:
	docker build . -t "ctownsend/franken-mon"
 
	docker run --rm -v "$(PWD)":/franken-mon ctownsend/franken-mon sh -c \
	"cd /franken-mon && go build -o franken-mon-linux"

	docker run --rm -v "$(PWD)":/franken-mon --rm -e GOOS=darwin -e GOARCH=amd64 ctownsend/franken-mon sh -c \
	"cd /franken-mon && go build -o franken-mon-darwin"

clean:
	rm -f franken-mon-linux && rm -f franken-mon-darwin
	docker rmi ctownsend/franken-mon



