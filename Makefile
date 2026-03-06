.PHONY: dev build install clean test test-v test-update demo demo-quick screenshots

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
# Requires: vhs, ffmpeg, google-chrome
# Install vhs: go install github.com/charmbracelet/vhs@latest
# Install ffmpeg: sudo apt install ffmpeg

demo: build
	mkdir -p demos/screenshots
	PATH="$(PWD)/dist:$(PATH)" vhs demos/demo.tape

demo-quick: build
	PATH="$(PWD)/dist:$(PATH)" vhs demos/quick.tape

screenshots: build
	mkdir -p demos/screenshots
	PATH="$(PWD)/dist:$(PATH)" vhs demos/screenshots.tape
