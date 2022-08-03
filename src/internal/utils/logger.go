package utils

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Loggers struct {
	Logger  zerolog.Logger
	Bot     *tgbotapi.BotAPI
	adminId int64
}

func NewLogger(bot *tgbotapi.BotAPI, adminId int64) Loggers {
	writers := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: time.RFC1123,
	})
	Logger := zerolog.New(writers)
	Logger.Level(zerolog.InfoLevel)
	return Loggers{
		Logger:  Logger,
		Bot:     bot,
		adminId: adminId,
	}
}

func (l Loggers) SendAlertToAdmin(prefix string, err error) {
	msg := tgbotapi.NewMessage(l.adminId, fmt.Sprintf("Error %s %+v", prefix, err))
	if _, err := l.Bot.Send(msg); err != nil {
		l.Logger.Error().Timestamp().Err(err).Msg("SendAlertToAdmin error")
		return
	}
	return
}
