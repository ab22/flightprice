SHELL := /bin/bash
-include .env
export

.PHONY: up
up:
	docker compose --progress=plain -f ./build/local/docker-compose.yml up --remove-orphans

.PHONY: down
down:
	docker compose -f ./build/local/docker-compose.yml down --remove-orphans

.PHONY: build
build:
	docker compose -f ./build/local/docker-compose.yml build

.PHONY: flogs
flogs:
	docker compose -f ./build/local/docker-compose.yml logs -f

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: ping
ping:
	curl -i http://localhost:8080/ping

.PHONY: login
login:
	curl -i http://localhost:8080/login

.PHONY: flights
flights:
	curl -i --header "X-API-Token: ${API_TOKEN}" http://localhost:8080/flights/search

.PHONY: flights-invalid
flights-invalid:
	curl -i --header "X-API-Token: invalidtoken" http://localhost:8080/flights/search

.PHONY: ws
ws:
	curl --include \
	     --no-buffer \
	     --header "Connection: Upgrade" \
	     --header "Upgrade: websocket" \
	     --header "Host: localhost:8080" \
	     --header "Origin: http://localhost:8080" \
	     --header "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
	     --header "Sec-WebSocket-Version: 13" \
		"http://localhost:8080/subscribe/3"
