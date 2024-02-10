package dbwrapper

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mderler/simple-go-backend/internal/db"
	"go.uber.org/fx"
)

func NewPostgresQueries(lc fx.Lifecycle) *db.Queries {
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

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return err
		},
		OnStop: func(ctx context.Context) error {
			conn.Close()
			return nil
		},
	})

	return db.New(conn)
}
