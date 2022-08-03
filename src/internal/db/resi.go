package db

import (
	"context"
	"fmt"
	"src/internal/entity"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

var (
	querySearchResi          = "SELECT user_id, no_resi, vendor, history, chat_id, completed, name from resi WHERE completed = false and is_deleted = false "
	querySearchAllResi       = "SELECT user_id, no_resi, vendor, history, chat_id, completed, name from resi "
	queryInsertResi          = "INSERT INTO resi (user_id, no_resi, vendor, chat_id, name) VALUES ($1, $2, $3, $4, $5)"
	queryUpdateResi          = "UPDATE resi SET history = $1, completed = $2, updated_at = NOW() WHERE no_resi = $3 and user_id = $4"
	queryUpdateResiCompleted = "UPDATE resi SET completed = $1, updated_at = NOW() WHERE no_resi = $2 and user_id = $3"
	queryDeleteResi          = "UPDATE resi SET is_deleted = true WHERE user_id = $1 AND no_resi = $2"
)

func (db Database) GetDataResi(ctx context.Context) ([]*entity.Resi, error) {
	var allResi []*entity.Resi
	err := pgxscan.Select(ctx, db.database, &allResi, querySearchResi)
	if err != nil {
		db.Logger.Logger.Error().Timestamp().Err(err).Msg("GetDataResi failed")
		return nil, err
	}

	return allResi, nil
}

func (db Database) GetDataResiByUserId(ctx context.Context, userId string) ([]*entity.Resi, error) {
	query := querySearchAllResi + "where is_deleted = false and user_id = $1 and updated_at >= $2"
	tanggal := time.Now().Add(-time.Hour * 24 * 50)
	var allResi []*entity.Resi
	err := pgxscan.Select(ctx, db.database, &allResi, query, userId, tanggal)
	if err != nil {
		db.Logger.Logger.Error().Timestamp().Err(err).Msg("GetDataResiByUserId failed")
		return nil, err
	}

	return allResi, nil
}

func (db Database) UpdateDataResi(ctx context.Context, resi entity.Resi) error {
	db.Logger.Logger.Info().Timestamp().Msgf("UpdateDataResi %v", resi)
	what, err := db.database.Exec(ctx, queryUpdateResi, resi.History, resi.Completed, resi.NoResi, resi.UserID)
	if err != nil {
		db.Logger.Logger.Error().Timestamp().Err(err).Msg("Update Resi failed")
	}
	db.Logger.Logger.Info().Timestamp().Msgf("What %v", what)

	return err
}

func (db Database) UpdateDataResiCompleted(ctx context.Context, completed bool, resi, userId string) error {
	db.Logger.Logger.Info().Timestamp().Msgf("UpdateDataResi %v", resi)
	what, err := db.database.Exec(ctx, queryUpdateResiCompleted, completed, resi, userId)
	if err != nil {
		db.Logger.Logger.Error().Timestamp().Err(err).Msg("Update Resi Completed failed")
	}
	db.Logger.Logger.Info().Timestamp().Msgf("What %v", what)

	return err
}

func (db Database) SaveDataResi(ctx context.Context, user entity.User, vendor, noResi, chatId, name string) error {
	db.Logger.Logger.Info().Timestamp().Msg("Save Data Resi")
	_, err := db.database.Exec(ctx, queryInsertResi, user.UserID, noResi, vendor, chatId, name)
	if err != nil {
		db.Logger.Logger.Error().Timestamp().Err(err).Msg("Save Resi failed")
	}
	return err
}

func (db Database) DeleteDataResi(ctx context.Context, userId string) error {
	query := fmt.Sprintf("%sWHERE user_id = $1 ORDER BY created_at DESC LIMIT 1", querySearchAllResi)
	res := entity.Resi{}
	// user_id, no_resi, vendor, history, chat_id, completed, name
	err := db.database.QueryRow(ctx, query, userId).
		Scan(&res.UserID, &res.NoResi, &res.Vendor, &res.History, &res.ChatID, &res.Completed, &res.Name)
	if err != nil {
		return err
	}
	_, err = db.database.Exec(ctx, queryDeleteResi, res.UserID, res.NoResi)
	return err
}

func (db Database) DeleteDataResiByResi(ctx context.Context, userId, resi string) error {
	_, err := db.database.Exec(ctx, queryDeleteResi, userId, resi)
	return err
}
