package test_helpers

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/stretchr/testify/mock"
)

// mockGameResultDAO is a mock type for the GameResultDAO type
type mockGameResultDAO struct {
	mock.Mock
}

// NewMockGameResultDAO creates a new instance of MockManager
func NewMockGameResultDAO() *mockGameResultDAO {
	return &mockGameResultDAO{}
}

func (m *mockGameResultDAO) CreateGameResult(
	ctx context.Context,
	userId uuid.UUID,
	gameStatus entity.GameStatus,
	amount float64,
	transactionSource entity.TransactionSource,
	transactionID string) (*entity.GameResult, error) {

	args := m.Called(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	if arg := args.Get(0); arg != nil {
		return arg.(*entity.GameResult), nil
	}
	return nil, args.Error(1)
}

func (m *mockGameResultDAO) ValidateGameResults(ctx context.Context, totalGamesToCancel int) error {
	args := m.Called(ctx, totalGamesToCancel)
	return args.Error(0)
}
