package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID    `db:"id"`
	Email                string       `db:"email"`
	Balance              float64      `db:"balance"`
	LastGameResultAt     sql.NullTime `db:"last_game_result_at"`
	GamesResultValidated sql.NullBool `db:"games_result_validated"`
	CreatedAt            time.Time    `db:"created_at"`
}
