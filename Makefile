APP_NAME=structgen

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)/
