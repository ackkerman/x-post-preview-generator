.PHONY: all build test fmt vet lint tidy sample

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
		-icon "https://pbs.twimg.com/profile_images/1661201415899951105/azNjKOSH_400x400.jpg" \
		-date "5:50 AM - Mar 22, 2006" \
		-cta "Read 16K replies" \
		-output "samples/jack.png"
