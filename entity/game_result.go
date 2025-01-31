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

	if source != TransactionSourceGame &&
		source != TransactionSourceServer &&
		source != TransactionSourcePayment {
		return nil
	}
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
	ID                int               `db:"id"`
	UserID            uuid.UUID         `db:"user_id"`
	GameStatus        GameStatus        `db:"game_status"`
	ValidationStatus  ValidationStatus  `db:"validation_status"`
	TransactionSource TransactionSource `db:"transaction_source"`
	TransactionID     string            `db:"transaction_id"`
	Amount            float64           `db:"amount" `
	CreatedAt         time.Time         `db:"created_at"`
}

func (dm *GameResult) ShouldBeCanceled() bool {

	// Check it the ID is odd
	return dm.ID%2 != 0
}
