package utils

import (
	"fmt"
	"src/internal/entity"

	"github.com/aprilahijriyan/gojne/lib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

func CreateMessageSicepat(resi entity.Resi, his []entity.TrackHistory) string {
	message := fmt.Sprintf("No Resi: %s \nKurir: %s\nNama: %s\n\n", resi.NoResi, resi.Vendor, resi.Name)

	for _, h := range his {
		middle := h.City
		if h.City == "" {
			middle = h.ReceiverName
		}
		message += fmt.Sprintf("%s-%s \n%s\n\n", h.Status, middle, h.Date)
	}
	return message
}

func CreateMessageJNE(resi entity.Resi, res lib.DetailTracking) string {
	message := fmt.Sprintf("No Resi: %s\nNama: %s\n", resi.NoResi, resi.Name)

	for _, del := range res.Data.Delivery {
		message += fmt.Sprintf("%s: %s\n", del.Title, del.Value)
	}

	message += "\n\n"

	for _, h := range res.Data.History.Data {
		message += fmt.Sprintf("%s\n%s\n\n", h.Title, h.Date)
	}
	return message
}

func SendUpdateResiToUser(bot *tgbotapi.BotAPI, chatID int64, text string, log zerolog.Logger) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Error().Timestamp().Err(err).Msg("SendMessage error")
		return err
	}
	return nil
}

func CreateMessageGetList(resi []*entity.Resi) string {
	message := fmt.Sprintf("List Resi\n\n")

	for i, r := range resi {
		message += fmt.Sprintf("%d.\nNo Resi: %s\nNama: %s\nStatus: ", i+1, r.NoResi, r.Name)
		if r.Completed {
			message += "Tiba"
		} else {
			message += "Dalam Proses"
		}
		message += "\n\n"
	}
	return message
}
