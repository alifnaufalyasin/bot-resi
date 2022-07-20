package db

import (
	"context"
	"fmt"
	"src/internal/entity"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

var (
	querySearchUser = "SELECT id, user_id, username, name from users"
	queryInsertUser = "INSERT INTO users (user_id, username, name) VALUES ($1, $2, $3)"
	queryUpdateUser = "UPDATE users SET username = $1, name = $2, updated_at = NOW() WHERE user_id = $3"
)

func (db Database) SaveUser(ctx context.Context, user *tgbotapi.User) (entity.User, error) {
	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	u := entity.User{
		UserID:   strconv.FormatInt(user.ID, 10),
		Username: user.UserName,
		Name:     name,
	}
	err := db.database.QueryRow(ctx, fmt.Sprintf("%s where user_id = $1", querySearchUser), strconv.FormatInt(user.ID, 10)).
		Scan(&u.Id, &u.UserID, &u.Username, &u.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = db.database.Exec(ctx, queryInsertUser, u.UserID, u.Username, u.Name)
			if err != nil {
				db.logger.Error().Timestamp().Err(err).Msg("Insert user failed")
			}
		} else {
			db.logger.Error().Timestamp().Err(err).Msg("QueryRowUser failed")
		}
	}
	change := false
	if strconv.FormatInt(user.ID, 10) != u.UserID {
		change = true
		u.UserID = strconv.FormatInt(user.ID, 10)
	}
	if user.UserName != u.Username {
		change = true
		u.Username = user.UserName
	}

	if name != u.Name {
		change = true
		u.Name = name
	}
	if change {
		_, err = db.database.Exec(ctx, queryUpdateUser, u.Username, u.Name, u.UserID)
		if err != nil {
			db.logger.Error().Timestamp().Err(err).Msg("Update user failed")
		}
	}

	return u, nil
}
