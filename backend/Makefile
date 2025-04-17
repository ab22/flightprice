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

.PHONY: lint
	golangci-lint run ./...
