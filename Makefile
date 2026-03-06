.PHONY: dev build install clean test test-v test-update demo demo-quick screenshots playground

dev:
	go run ./cmd/vs

build:
	go build -o dist/vs ./cmd/vs

install:
	go install ./cmd/vs

clean:
	rm -rf dist/ .playground/

# --- Tests (teatest) ---

test:
	go test ./... -timeout 30s

test-v:
	go test ./... -v -timeout 30s

test-update:
	go test ./... -update -timeout 30s

# --- Playground (teste manual do fluxo completo) ---
# Cria .playground/ limpo, reseta onboarding, roda vs
playground: build
	bash scripts/playground.sh

# --- Demos (vhs) ---
# Requires: vhs, ffmpeg
# Install vhs: go install github.com/charmbracelet/vhs@latest

demo: build
	mkdir -p demos/screenshots
	PATH="$(PWD)/dist:$(HOME)/bin:$(PATH)" vhs demos/demo.tape

demo-quick: build
	PATH="$(PWD)/dist:$(HOME)/bin:$(PATH)" vhs demos/quick.tape

screenshots: build
	mkdir -p demos/screenshots
	PATH="$(PWD)/dist:$(HOME)/bin:$(PATH)" vhs demos/screenshots.tape
	PATH="$(PWD)/dist:$(HOME)/bin:$(PATH)" vhs demos/onboarding.tape
