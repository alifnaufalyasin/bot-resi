package controller

import (
	"encoding/json"
	"src/internal/entity"
	"src/internal/utils"

	"github.com/rs/zerolog"
)

func UpdateResiSicepat(res *entity.Resi, uri string, log zerolog.Logger) (msg string, send bool, err error) {
	send = false
	resi, raw, err := entity.GetResiSicepatHistory(uri, res.NoResi)
	if err != nil {
		log.Error().Timestamp().Err(err).Msg("GetResiHistory error")
		return
	}

	if res.History != "" {
		resiDb := entity.SicepatRes{}
		err = json.Unmarshal([]byte(res.History), &resiDb)
		if err != nil {
			log.Error().Timestamp().Err(err).Msg("Unmarshal error")
		}

		if len(resi.Sicepat.Result.TrackHistory) > len(resiDb.Sicepat.Result.TrackHistory) {
			send = true
			res.History = raw

			if resi.Sicepat.Result.TrackHistory[len(resi.Sicepat.Result.TrackHistory)-1].Status == "DELIVERED" {
				res.Completed = true
			}
		}
	} else {
		send = true
		if len(resi.Sicepat.Result.TrackHistory) > 0 {
			res.History = raw
		}
	}

	msg = utils.CreateMessageSicepat(*res, resi.Sicepat.Result.TrackHistory)
	return
}
