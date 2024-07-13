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
	Amount float64,
	transactionID string) (*entity.GameResult, error) {

	args := m.Called(ctx, userId, gameStatus, Amount, transactionID)

	if len(args) > 0 && args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}

	return args.Get(0).(*entity.GameResult), nil
}

func (m *mockGameResultDAO) ValidateGameResults(ctx context.Context, totalGamesToCancel int) error {
	args := m.Called(ctx, totalGamesToCancel)
	return args.Error(0)
}
