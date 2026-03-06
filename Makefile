.PHONY: dev build install clean test test-v demo demo-quick

dev:
	go run ./cmd/vs

build:
	go build -o dist/vs ./cmd/vs

install:
	go install ./cmd/vs

clean:
	rm -rf dist/

# --- Tests (teatest) ---

test:
	go test ./... -timeout 30s

test-v:
	go test ./... -v -timeout 30s

test-update:
	go test ./... -update -timeout 30s

# --- Demos (vhs) ---

demo: build
	mkdir -p demos/screenshots
	vhs demos/demo.tape

demo-quick: build
	vhs demos/quick.tape
