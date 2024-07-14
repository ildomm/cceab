package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/jmoiron/sqlx"
)

type Querier interface {
	Close()
	WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) (err error)

	InsertGameResult(ctx context.Context, txn sqlx.Tx, gameResult entity.GameResult) (int, error)

	LockUserRow(ctx context.Context, txn sqlx.Tx, userId uuid.UUID) error
	SelectUser(ctx context.Context, userId uuid.UUID) (*entity.User, error)
	SelectUsersByValidationStatus(ctx context.Context, validationStatus bool) ([]entity.User, error)

	CheckTransactionID(ctx context.Context, transactionId string) (bool, error)
	SelectGameResultsByUser(ctx context.Context, userId uuid.UUID, validationStatus entity.ValidationStatus) ([]entity.GameResult, error)

	UpdateUserBalance(ctx context.Context, txn sqlx.Tx, userId uuid.UUID, balance float64, validationStatus bool) error
	UpdateGameResult(ctx context.Context, txn sqlx.Tx, gameResultId int, validationStatus entity.ValidationStatus) error
}
