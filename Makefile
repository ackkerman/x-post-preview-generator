.PHONY: all build test fmt vet lint tidy sample ui-wasm

WASM_EXEC := $(shell go env GOROOT)/lib/wasm/wasm_exec.js
WASM_EXEC_FALLBACK := $(shell go env GOROOT)/misc/wasm/wasm_exec.js

all: fmt lint test build

build:
	go build -o xpostgen ./cmd/xpostgen

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

lint: vet

tidy:
	go mod tidy

sample:
	mkdir -p samples
	go run ./cmd/xpostgen \
		-text "just setting up my twttr" \
		-name "jack" \
		-id "jack" \
		-verified \
		-like-count "262K" \
		-width-mode tight \
		-icon "https://pbs.twimg.com/profile_images/1661201415899951105/azNjKOSH_400x400.jpg" \
		-date "5:50 AM - Mar 22, 2006" \
		-cta "Read 16K replies" \
		-output "samples/jack.png"

	go run ./cmd/xpostgen \
		-text "just setting up my twttr" \
		-name "jack" \
		-id "jack" \
		-verified \
		-like-count "262K" \
		-width-mode tight \
		-icon "https://pbs.twimg.com/profile_images/1661201415899951105/azNjKOSH_400x400.jpg" \
		-date "5:50 AM - Mar 22, 2006" \
		-cta "Read 16K replies" \
		-output "samples/jack.svg"

	go run ./cmd/xpostgen \
		-text "just setting up my twttr" \
		-name "jack" \
		-id "jack" \
		-verified \
		-simple \
		-width-mode tight \
		-icon "https://pbs.twimg.com/profile_images/1661201415899951105/azNjKOSH_400x400.jpg" \
		-date "5:50 AM - Mar 22, 2006" \
		-cta "Read 16K replies" \
		-output "samples/simple-jack.svg"

ui-wasm:
	mkdir -p ui/public/wasm
	GOOS=js GOARCH=wasm go build -o ui/public/wasm/xpostgen.wasm ./cmd/xpostgen-wasm
	if [ -f "$(WASM_EXEC)" ]; then \
		cp "$(WASM_EXEC)" ui/public/wasm/wasm_exec.js; \
	else \
		cp "$(WASM_EXEC_FALLBACK)" ui/public/wasm/wasm_exec.js; \
	fi

ui-install:
	pnpm --prefix ui install

ui-dev: ui-install ui-wasm
	pnpm --prefix ui dev --port 3001

ui-build: ui-install ui-wasm
	pnpm --prefix ui install