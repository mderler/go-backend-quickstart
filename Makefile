-include .env
export

.PHONY: create-migration apply-migrations run-postgres start-postgres stop-postgres sqlc-gen

GOOSE=github.com/pressly/goose/v3/cmd/goose@latest

DOCKER_CLIENT=docker

create-migration:
	go run $(GOOSE) -dir ./migrations create $(name) sql

apply-migrations:
	go run $(GOOSE) -dir ./migrations postgres \
		"host=$(POSTGRES_HOST) \
		user=$(POSTGRES_USER) \
		password=$(POSTGRES_PASSWORD) \
		dbname=$(POSTGRES_DB) sslmode=disable" up

run-postgres:
	$(DOCKER_CLIENT) run --name postgres \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p 5432:5432 -d postgres:alpine

start-postgres:
	$(DOCKER_CLIENT) start postgres

stop-postgres:
	$(DOCKER_CLIENT) stop postgres
	$(DOCKER_CLIENT) rm postgres

sqlc-gen:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

swag-init:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go