package db

import (
	"context"
	"os"
	"src/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Database struct {
	database *pgxpool.Pool
	logger   zerolog.Logger
}

func InitDB(ctx context.Context, cfg config.Config, log zerolog.Logger) Database {
	conn, err := pgxpool.Connect(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Msgf("%v\nUnable to connect to database: %v\n", os.Stderr, err)
	}
	return Database{
		database: conn,
		logger:   log,
	}
}
