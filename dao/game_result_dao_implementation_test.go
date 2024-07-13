package dao

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/ildomm/cceab/entity"
	"github.com/ildomm/cceab/test_helpers"
)

func TestCreateGameResultSuccess(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	mockGameResultDAO := test_helpers.NewMockGameResultDAO()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()
	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame
	transactionID := "unique-transaction-id"

	// Mock successful interactions
	mockQuerier.On("CheckTransactionID", ctx, transactionID).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return(&entity.User{
		ID:      userId,
		Balance: 200.0,
	}, nil)
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error"))
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userId).Return(nil)
	mockQuerier.On("InsertGameResult", ctx, mock.Anything, mock.Anything).Return(uuid.New(), nil)
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userId, mock.Anything, false).Return(nil)

	_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	assert.NoError(t, err, "CreateGameResult should not return an error")
	mockQuerier.AssertExpectations(t)
	mockGameResultDAO.AssertExpectations(t)

	// Count the game results
	mockQuerier.On("SelectGameResultsByUser", ctx, mock.Anything, mock.Anything)
	games, err := mockQuerier.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusPending)
	assert.NoError(t, err)
	assert.Len(t, games, 1, "There should be one game result")
}

func TestCreateGameResultTransactionIDExists(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	mockGameResultDAO := test_helpers.NewMockGameResultDAO()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()
	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame
	transactionID := "existing-transaction-id"

	// Mock transaction ID already exists
	mockQuerier.On("CheckTransactionID", ctx, transactionID).Return(true, nil)

	_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	assert.EqualError(t, err, entity.ErrTransactionIdExists.Error(), "CreateGameResult should return ErrTransactionIdExists")
	mockQuerier.AssertExpectations(t)
	mockGameResultDAO.AssertExpectations(t)
}

func TestCreateGameResultUserNotFound(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	mockGameResultDAO := test_helpers.NewMockGameResultDAO()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()
	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame
	transactionID := "unique-transaction-id"

	// Mock user not found
	mockQuerier.On("CheckTransactionID", ctx, transactionID).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return(nil, entity.ErrUserNotFound)

	_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	assert.EqualError(t, err, entity.ErrUserNotFound.Error(), "CreateGameResult should return ErrUserNotFound")
	mockQuerier.AssertExpectations(t)
	mockGameResultDAO.AssertExpectations(t)
}

func TestCreateGameResultInsufficientBalance(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	mockGameResultDAO := test_helpers.NewMockGameResultDAO()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()
	userId := uuid.New()
	gameStatus := entity.GameStatusLost // Assuming this triggers the balance check
	amount := 300.0                     // Assuming the user's balance is less than this amount
	transactionSource := entity.TransactionSourceGame
	transactionID := "unique-transaction-id"

	// Mock user with insufficient balance
	mockQuerier.On("CheckTransactionID", ctx, transactionID).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return(&entity.User{
		ID:      userId,
		Balance: 200.0,
	}, nil)

	_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	assert.EqualError(t, err, entity.ErrUserNegativeBalance.Error(), "CreateGameResult should return ErrUserNegativeBalance")
	mockQuerier.AssertExpectations(t)
	mockGameResultDAO.AssertExpectations(t)
}

func TestCreateGameResultDatabaseError(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	mockGameResultDAO := test_helpers.NewMockGameResultDAO()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()
	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame
	transactionID := "unique-transaction-id"

	// Mock successful interactions except for InsertGameResult
	mockQuerier.On("CheckTransactionID", ctx, transactionID).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return(&entity.User{
		ID:      userId,
		Balance: 200.0,
	}, nil)
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error"))
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userId).Return(nil)
	mockQuerier.On("InsertGameResult", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)

	assert.EqualError(t, err, entity.ErrCreatingGameResult.Error(), "CreateGameResult should return ErrCreatingGameResult")
	mockQuerier.AssertExpectations(t)
	mockGameResultDAO.AssertExpectations(t)

	// Count the game results
	mockQuerier.On("SelectGameResultsByUser", ctx, mock.Anything, mock.Anything)
	games, err := mockQuerier.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusPending)
	assert.NoError(t, err)
	assert.Len(t, games, 0, "There should be none game result")
}
