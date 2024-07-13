package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
)

type Querier interface {
	Close()

	InsertGameResult(ctx context.Context, gameResult entity.GameResult) error

	SelectUser(ctx context.Context, userId uuid.UUID) (*entity.User, error)
	SelectUsersByValidationStatus(ctx context.Context, validationStatus bool) ([]entity.User, error)

	GameResultExists(ctx context.Context, transactionId string) (bool, error)
	// SelectGameResultByTransaction(ctx context.Context, transactionId string) (*entity.GameResult, error)
	SelectGameResultsByUser(ctx context.Context, userId uuid.UUID, validationStatus entity.ValidationStatus) ([]entity.GameResult, error)

	UpdateUserBalance(ctx context.Context, userId uuid.UUID, amount float64, validationStatus bool) error
	UpdateGameResult(ctx context.Context, gameResultId uuid.UUID, validationStatus entity.ValidationStatus) error
}
