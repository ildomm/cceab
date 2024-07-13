package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
)

type GameResultDAO interface {
	CreateGameResult(ctx context.Context, userId uuid.UUID, gameStatus entity.GameStatus, amount float64, transactionSource entity.TransactionSource, transactionID string) (*entity.GameResult, error)
	ValidateGameResults(ctx context.Context, totalGamesToCancel int) error
}
