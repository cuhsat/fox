bin = /usr/local/bin/

all: build install symlink

build:
	mkdir -p ./bin
	go build -v -race -o ./bin/cu cmd/cu/main.go

install: build
	sudo mkdir -p $(bin)
	sudo cp ./bin/cu $(bin)

symlink: install
	sudo ln -s $(bin)/cu $(bin)/icu

clean:
	rm -rf ./bin
	go clean

.PHONY: all clean
