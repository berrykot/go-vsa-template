wire:
	wire ./cmd/app

build: wire
	go build -o bin/app ./cmd/app

run: wire
	go run ./cmd/app

tidy:
	go mod tidy

all: tidy wire build