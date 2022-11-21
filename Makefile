

.PHONY: build
all: build

build:
	@echo "build"
	go build -o build/gnobot ./cmd/gnodiscord

build_linux:
	@echo "build"
	GOOS=linux GOARCH=amd64 go build -o build/gnobot_linux ./cmd/gnodiscord
