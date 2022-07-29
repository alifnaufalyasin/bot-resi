package controller

import (
	"encoding/json"
	"src/internal/entity"
	"src/internal/utils"
	"strings"

	"github.com/aprilahijriyan/gojne/lib"
	"github.com/rs/zerolog"
)

func UpdateResiJNE(res *entity.Resi, log zerolog.Logger) (msg string, send bool, err error) {
	resi := entity.GetResiJNEHistory(res.NoResi)

	log.Info().Timestamp().Msgf("%+v", resi)

	send = false

	if res.History != "" {
		resiDb := lib.DetailTracking{}
		err = json.Unmarshal([]byte(res.History), &resiDb)
		if err != nil {
			log.Error().Timestamp().Err(err).Msg("Unmarshal error")
		}

		if len(resi.Data.History.Data) > len(resiDb.Data.History.Data) {
			send = true
			out, _ := json.Marshal(resi)
			res.History = string(out)

			if strings.Split(resi.Data.History.Data[len(resi.Data.History.Data)-1].Title, " ")[0] == "DELIVERED" {
				res.Completed = true
			}
		}
	} else {
		send = true
		if len(resi.Data.History.Data) > 0 {
			out, _ := json.Marshal(resi)
			res.History = string(out)
		}
	}

	msg = utils.CreateMessageJNE(*res, resi)

	return
}
