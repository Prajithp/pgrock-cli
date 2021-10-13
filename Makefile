.PHONY: build fmt all clean

fmt:
	go fmt .

build:
	go build

clean:
	go clean -i

all: fmt build