# go-backend-quickstart

## Introduction

The goal of this template is to create a solid starting point for a Go backend.

This template will give you a basic structure including an example todo-app that implements the following things:

- A lightweight router with [go-chi/chi](https://github.com/go-chi/chi/)
- Generation of swagger docs with [swaggo/swag](https://github.com/swaggo/swag) and [swaggo/http-swagger](https://github.com/swaggo/http-swagger)
- Request validation with [go-playground/validator](https://github.com/go-playground/validator)
- Migration of a [PostgreSQL](https://www.postgresql.org/) DB with [pressly/goose](https://github.com/pressly/goose)
- Generation of type-safe interfaces from SQL with [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc) with the [jackc/pgx](https://github.com/jackc/pgx) driver
- Management of environment variables with [jogo/godotenv](https://github.com/joho/godotenv)
- Automation of some tasks with a Makefile and [Docker](https://www.docker.com/)

### Disclaimer

This template is far from perfect. I am __not__ an experienced Go developer. I am open for suggestions!

## Usage

To run the example, follow these steps:

1. Run `go mod tidy`.
2. Make a copy of `.env.example` with the name `.env`.
3. Set the required environment variables.
4. Run `make run-postgres` (In the Makefile, set DOCKER_CLIENT to podman if you're using podman).
5. Run `make apply-migrations`.
6. Run `go run cmd/api/main.go`.

The documentation is now accessible at http://localhost:3000/swagger/.

### The Todo-Example

If you think that the example code is bloated then you are probably right...

...but I wanted to showcase as much as possible in terms of what is achievable
with this stack and how to implement it.

#### What it does

There are Users and Todos. A User can create a Todo and Users can be assigned
to Todos.

__Updates to the documentation on this templates will follow soon__
