package db

import (
	"context"
	"os"
	"src/internal/config"
	"src/internal/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	database *pgxpool.Pool
	Logger   utils.Loggers
}

func InitDB(ctx context.Context, cfg config.Config, log utils.Loggers) Database {
	conn, err := pgxpool.Connect(ctx, cfg.Database.URL)
	if err != nil {
		log.Logger.Fatal().Msgf("%v\nUnable to connect to database: %v\n", os.Stderr, err)
		log.SendAlertToAdmin("iditDB", err)
	}
	return Database{
		database: conn,
		Logger:   log,
	}
}
