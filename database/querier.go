package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/domain"
)

type Querier interface {
	Close()

	InsertGameResult(ctx context.Context, gameResult domain.GameResult) error

	SelectUser(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	SelectUsersByValidationStatus(ctx context.Context, validationStatus bool) ([]domain.User, error)

	SelectGameResultByTransaction(ctx context.Context, transactionId string) (*domain.GameResult, error)
	SelectGameResultsByUser(ctx context.Context, userId uuid.UUID, validationStatus domain.ValidationStatus) ([]domain.GameResult, error)

	UpdateUserBalance(ctx context.Context, userId uuid.UUID, amount float64, validationStatus bool) error
	UpdateGameResult(ctx context.Context, gameResultId uuid.UUID, validationStatus domain.ValidationStatus) error
}
