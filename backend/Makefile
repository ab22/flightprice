.PHONY: up
up:
	docker compose --progress=plain -f ./build/local/docker-compose.yml up -d --remove-orphans

.PHONY: down
down:
	docker compose down --remove-orphans -f ./build/local/docker-compose.yml

.PHONY: build
build:
	docker compose --progress=plain -f ./build/local/docker-compose.yml build --no-cache

.PHONY: lint
	golangci-lint run ./...
