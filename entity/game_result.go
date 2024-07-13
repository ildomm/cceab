package entity

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"time"
)

type GameStatus string
type ValidationStatus string
type TransactionSource string

const (
	GameStatusWin  GameStatus = "win"
	GameStatusLost GameStatus = "lost"
)

const (
	ValidationStatusPending  ValidationStatus = "pending"
	ValidationStatusAccepted ValidationStatus = "accepted"
	ValidationStatusCanceled ValidationStatus = "canceled"
)

const (
	TransactionSourceGame    TransactionSource = "game"
	TransactionSourceServer  TransactionSource = "server"
	TransactionSourcePayment TransactionSource = "payment"
)

func ParseTransactionSource(value interface{}) *TransactionSource {
	source := TransactionSource(value.(string))
	return &source
}

func (e *GameStatus) Scan(value interface{}) error {
	*e = GameStatus(value.(string))
	return nil
}

func (e GameStatus) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *ValidationStatus) Scan(value interface{}) error {
	*e = ValidationStatus(value.(string))
	return nil
}

func (e ValidationStatus) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *TransactionSource) Scan(value interface{}) error {
	*e = TransactionSource(value.(string))
	return nil
}

func (e TransactionSource) Value() (driver.Value, error) {
	return string(e), nil
}

type GameResult struct {
	ID                int               `db:"id" json:"id"`
	UserID            uuid.UUID         `db:"user_id" json:"user_id"`
	GameStatus        GameStatus        `db:"game_status" json:"game_status"`
	ValidationStatus  ValidationStatus  `db:"validation_status" json:"validation_status"`
	TransactionSource TransactionSource `db:"transaction_source" json:"transaction_source"`
	TransactionID     string            `db:"transaction_id" json:"transaction_id"`
	Amount            float64           `db:"amount" json:"amount"`
	CreatedAt         time.Time         `db:"created_at" json:"created_at"`
}
