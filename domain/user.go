package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID    `db:"id" json:"id"`
	Email                string       `db:"email" json:"email"`
	Balance              float64      `db:"balance" json:"balance"`
	LastGameResultAt     sql.NullTime `db:"last_game_result_at" json:"last_game_result_at"`
	GamesResultValidated sql.NullBool `db:"games_result_validated" json:"games_result_validated"`
	CreatedAt            time.Time    `db:"created_at" json:"created_at"`
}
