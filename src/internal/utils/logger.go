package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	writers := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: time.RFC1123,
	})
	Logger := zerolog.New(writers)
	Logger.Level(zerolog.InfoLevel)
	return Logger
}
