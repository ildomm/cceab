package dao

import (
	"context"
	"errors"
	"github.com/ildomm/cceab/domain"
)

var ErrGameResultNotFound = errors.New("game result not found")
var ErrTransactionIdExists = errors.New("transaction id already exists")
var ErrInvalidGameStatus = errors.New("invalid game status")
var ErrInvalidTransactionSource = errors.New("invalid transaction source")
var ErrInvalidTransactionStatus = errors.New("invalid transaction status")

type GameResultDAO interface {
	CreateGameResult(ctx context.Context, gameResult domain.GameResult) (*domain.GameResult, error)
	ValidateGameResults(ctx context.Context, totalGamesToCancel int) error
}
