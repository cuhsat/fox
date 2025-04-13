bin = /usr/local/bin/

all: build install

build:
	mkdir -p ./bin
	go build -v -race -o ./bin/fx cmd/fx/main.go

install: build
	sudo mkdir -p $(bin)
	sudo cp ./bin/fx $(bin)

clean:
	rm -rf ./bin
	go clean

.PHONY: all clean
