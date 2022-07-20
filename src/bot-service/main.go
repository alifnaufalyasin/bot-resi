package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"src/internal/config"
	"src/internal/controller"
	"src/internal/db"
	"src/internal/utils"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logger := utils.NewLogger()
	cfg := config.GetConfig()
	db := db.InitDB(context.Background(), cfg, logger)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		logger.Panic().Timestamp().Err(err)
	}

	bot.Debug = false

	logger.Info().Timestamp().Msg(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Cron Function
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(1).Minute().SingletonMode().Do(controller.CheckStatusResi, context.Background(), db, bot, cfg.UriSicepat, logger)
	if err != nil {
		logger.Error().Timestamp().Err(err).Msg("Cron error")
	}
	s.StartImmediately().StartAsync()

	for update := range updates {
		u, err := db.SaveUser(context.Background(), update.SentFrom())
		if err != nil {
			logger.Error().Timestamp().Err(err).Msg("failed to save user")
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		MessageTelegram := update.Message.Text

		defer func() {
			if errr := recover(); errr != nil {
				msg.Text = "No resi JNE invalid"
				logger.Error().Fields(errr).Timestamp().Msg(msg.Text)
				if _, err := bot.Send(msg); err != nil {
					logger.Error().Err(err).Timestamp().Msg("error send")
				}
				err = db.DeleteDataResi(context.Background(), u.UserID)
				if err != nil {
					logger.Error().Err(err).Timestamp().Msg("Failed to delete invalid resi")
				}
			}
		}()

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "You can tracking your package with me \n\nkurir sicepat \nuse /sicepat <resi_no> <nama_paket>\n\nkurir jne\nuse /jne <resi_no> <nama_paket>\n\nuse /list to get list your tracking \n\nfor stop get update\nuse /stop <resi_no>\n\nfor resume get update \nuse /resume <resi_no>"
		case "status":
			msg.Text = "I'm ok."
		case "sicepat":
			res, err := controller.SaveDataResi(db, u, bot, MessageTelegram, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error Save Resi")
			}
			message, send, err := controller.UpdateResiSicepat(&res, cfg.UriSicepat, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error update resi")
			}
			err = controller.SendFirstResiUpdate(db, res, send, message, bot, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error send")
			}
			continue
		case "jne":
			res, err := controller.SaveDataResi(db, u, bot, MessageTelegram, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error Save Resi")
			}
			message, send, err := controller.UpdateResiJNE(&res, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error update resi")
			}
			err = controller.SendFirstResiUpdate(db, res, send, message, bot, logger)
			if err != nil {
				logger.Error().Err(err).Timestamp().Msg("error send")
			}
			continue
		case "stop":
			resi := strings.Split(MessageTelegram, " ")[1]
			err = db.UpdateDataResiCompleted(context.Background(), true, resi, strconv.FormatInt(update.Message.Chat.ID, 10))
			if err != nil {
				msg.Text = "Process stop gagal"
			} else {
				msg.Text = "Berhasil berhenti mendapatkan update resi. kirim /resume <no_resi> untuk melanjutkan mendapatkan update resi anda"
			}
			if _, err := bot.Send(msg); err != nil {
				logger.Panic().Err(err).Timestamp().Msg("error send")
			}
			continue
		case "resume":
			resi := strings.Split(MessageTelegram, " ")[1]
			err = db.UpdateDataResiCompleted(context.Background(), false, resi, strconv.FormatInt(update.Message.Chat.ID, 10))
			if err != nil {
				msg.Text = "Process resume gagal"
			} else {
				msg.Text = "Berhasil melanjutkan mendapatkan update resi. kirim /stop <no_resi> untuk berhenti mendapatkan update resi anda"
			}
			if _, err := bot.Send(msg); err != nil {
				logger.Panic().Err(err).Timestamp().Msg("error send")
			}
			continue
		case "list":
			resi, err := db.GetDataResiByUserId(context.Background(), u.UserID)
			if err != nil {
				msg.Text = "Get Resi gagal"
			}
			msg.Text = utils.CreateMessageGetList(resi)
			if _, err := bot.Send(msg); err != nil {
				logger.Panic().Err(err).Timestamp().Msg("error send list")
			}
			continue
		default:
			msg.Text = "I don't know that command"
		}
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				logger.Panic().Err(err).Timestamp().Msg("error send")
			}
		}
	}

}
