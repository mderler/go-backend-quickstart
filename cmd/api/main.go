package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mderler/simple-go-backend/internal/dbwrapper"
	"github.com/mderler/simple-go-backend/internal/handlers"
	"github.com/mderler/simple-go-backend/internal/validation"
	"go.uber.org/fx"
)

func newChiServer(lc fx.Lifecycle, userHandler *handlers.UserHandler, todoHandler *handlers.TodoHandler) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", userHandler)
		r.Mount("/todo", todoHandler)
	})

	srv := &http.Server{Addr: ":3000", Handler: r}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil

		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	validation.Init()

	fx.New(
		fx.Provide(dbwrapper.NewPostgresQueries, handlers.NewUserHandler, handlers.NewTodoHandler),
		fx.Invoke(newChiServer),
	).Run()
}
