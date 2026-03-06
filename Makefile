.PHONY: dev build install clean

dev:
	go run ./cmd/vs

build:
	go build -o dist/vs ./cmd/vs

install:
	go install ./cmd/vs

clean:
	rm -rf dist/
