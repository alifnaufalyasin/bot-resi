package controller

import (
	"context"
	"src/internal/db"
	"src/internal/entity"
	"src/internal/utils"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

func CheckStatusResi(ctx context.Context, db db.Database, bot *tgbotapi.BotAPI, uri string, log zerolog.Logger) {
	log.Info().Timestamp().Msg("Run Check")
	dataResi, err := db.GetDataResi(ctx)
	if err != nil {
		log.Error().Timestamp().Err(err).Msg("GetDataResi error")
	}
	if len(dataResi) == 0 {
		return
	}
	for _, r := range dataResi {
		res := *r
		msg := ""
		send := false
		if res.Vendor == "sicepat" {
			msg, send, err = UpdateResiSicepat(&res, uri, log)
		}

		if res.Vendor == "jne" {
			msg, send, err = UpdateResiJNE(&res, log)
		}

		if send {
			err = db.UpdateDataResi(ctx, res)
			if err != nil {
				log.Error().Timestamp().Err(err).Msg("UpdateDataResi error")
			}
			// Send message

			chatId, err := strconv.ParseInt(res.ChatID, 10, 64)
			if err != nil {
				log.Error().Timestamp().Err(err).Msg("ParseInt error")
			}
			utils.SendUpdateResiToUser(bot, chatId, msg, log)
		}
	}

	return

}

func SaveDataResi(db db.Database, u entity.User, bot *tgbotapi.BotAPI, Text string, log zerolog.Logger) (res entity.Resi, err error) {
	resi := strings.Split(Text, " ")[1]
	kurir := strings.Split(Text, " ")[0][1:]
	log.Info().Timestamp().Msg(kurir + " " + resi)
	chatId, _ := strconv.ParseInt(u.UserID, 10, 64)
	msg := tgbotapi.NewMessage(chatId, "")
	name := ""
	if len(strings.Split(Text, " ")) > 2 {
		name = strings.Join(strings.Split(Text, " ")[2:], " ")
	}
	err = db.SaveDataResi(context.Background(), u, kurir, resi, u.UserID, name)
	if err != nil {
		msg.Text = "Anda sudah pernah memasukkan no resi tersebut"
	} else {
		msg.Text = "Berhasil memasukkan no resi"
	}
	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Timestamp().Msg("error send")
	}

	res = entity.Resi{
		UserID:    u.UserID,
		NoResi:    resi,
		Vendor:    kurir,
		History:   "",
		ChatID:    u.UserID,
		Completed: false,
		Name:      name,
	}
	return
}

func SendFirstResiUpdate(db db.Database, res entity.Resi, send bool, message string, bot *tgbotapi.BotAPI, log zerolog.Logger) (err error) {
	if send {
		err = db.UpdateDataResi(context.Background(), res)
		if err != nil {
			log.Error().Timestamp().Err(err).Msg("UpdateDataResi error")
		}

		chatId, err := strconv.ParseInt(res.ChatID, 10, 64)
		if err != nil {
			log.Error().Timestamp().Err(err).Msg("ParseInt error")
		}
		err = utils.SendUpdateResiToUser(bot, chatId, message, log)
	}
	return
}
