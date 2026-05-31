.PHONY: build run test integration load lint

COMPOSE ?= docker compose
TOOLS   := $(COMPOSE) --profile tools run --rm tools

build:
	$(COMPOSE) build goboxd

run:
	$(COMPOSE) up goboxd

test:
	$(TOOLS) go test ./...

integration:
	$(TOOLS) go test -tags=integration ./tests/integration/...

load:
	$(TOOLS) go test -tags=load ./tests/load/...

lint:
	$(TOOLS) go vet ./...
	$(TOOLS) golangci-lint run ./...
