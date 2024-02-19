package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mderler/simple-go-backend/internal/db"
	"github.com/mderler/simple-go-backend/internal/handlers"

	_ "github.com/mderler/simple-go-backend/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Go Example API
// @version 1.0
// @description This is a sample API Server.
// @license.name MIT
// @BasePath /v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgxpool.New(
		context.TODO(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			"5432",
			os.Getenv("POSTGRES_DB"),
		),
	)
	if err != nil {
		log.Fatal("Error connecting with database")
	}

	defer conn.Close()

	queries := db.New(conn)

	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", handlers.NewUserHandler(queries))
		r.Mount("/todo", handlers.NewTodoHandler(queries))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))

	http.ListenAndServe(":3000", r)
}
