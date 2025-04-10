.PHONY: all clean

all: build install

build:
	mkdir -p ./bin
	go build -v -race -o ./bin/cu cmd/cu/main.go

install: build
	sudo mkdir -p /usr/local/bin/
	sudo cp ./bin/cu /usr/local/bin/

remove:
	sudo rm /usr/local/bin/cu

clean:
	rm -rf ./bin
	go clean
