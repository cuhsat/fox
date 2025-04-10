.PHONY: all clean

all: build

build:
	mkdir -p ./bin
	go build -v -race -o ./bin/cu cmd/cu/main.go

install: build
	chmod +x ./bin/cu
	sudo cp ./bin/cu /usr/local/bin/cu

remove:
	sudo rm /usr/local/bin/cu

clean:
	rm -rf ./bin
	go clean
