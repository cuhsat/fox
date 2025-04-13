src = ./bin/fx
dst = /usr/local/bin/

all: build install

build:
	mkdir -p ./bin
	go build -v -race -o $(src) cmd/fx/main.go

install: build
	sudo mkdir -p $(dst)
	sudo cp $(src) $(dst)

clean:
	rm -rf ./bin
	go clean

.PHONY: all clean
