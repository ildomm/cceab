package dao

import (
	"context"
	"errors"
	"fmt"
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

func TestValidateGameResultsSuccess(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()

	userID := uuid.New()

	// Mock select users with pending validation
	mockQuerier.On("SelectUsersByValidationStatus", ctx, false).Return([]entity.User{
		{
			ID:      userID,
			Balance: 100.0,
		},
	}, nil)

	// Mock lock user row
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userID).Return(nil)

	// Mock select game results for the user
	mockQuerier.On("SelectGameResultsByUser", ctx, userID, entity.ValidationStatusPending).Return([]entity.GameResult{
		{
			ID:                1,
			UserID:            userID,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusPending,
			TransactionSource: entity.TransactionSourceGame,
			TransactionID:     "tx123",
			Amount:            50.0,
		},
	}, nil)

	// Mock update game result
	mockQuerier.On("UpdateGameResult", ctx, mock.Anything, mock.Anything, entity.ValidationStatusCanceled).Return(nil)

	// Mock update user balance
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userID, mock.Anything, true).Return(nil)

	// Mock transaction operations
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error")).Return(nil)

	err := instance.ValidateGameResults(ctx, 1)

	assert.NoError(t, err, "ValidateGameResults should not return an error")
	mockQuerier.AssertExpectations(t)
}

func TestValidateGameResultsSuccessOnMultipleEntries(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()

	userID := uuid.New()
	balanceExpectedAdjustment := 2100.0

	// Mock select users with pending validation
	mockQuerier.On("SelectUsersByValidationStatus", ctx, false).Return([]entity.User{
		{
			ID:      userID,
			Balance: 2500.0, // 50 * 50 = 2500
		},
	}, nil)

	// Mock lock user row
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userID).Return(nil)

	// Populate list of game results
	totalEntries := 50
	totalGamesToCancel := 10 // 50 * 10 = 500
	var games []entity.GameResult

	// Add one lost game
	game := entity.GameResult{
		ID:                1,
		UserID:            userID,
		GameStatus:        entity.GameStatusLost,
		ValidationStatus:  entity.ValidationStatusPending,
		TransactionSource: entity.TransactionSourceGame,
		TransactionID:     fmt.Sprintf("tx%d", 1),
		Amount:            50.0,
	}
	games = append(games, game)

	// Add many wins
	for i := 0; i < totalEntries; i++ {
		game := entity.GameResult{
			ID:                i + 2,
			UserID:            userID,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusPending,
			TransactionSource: entity.TransactionSourceGame,
			TransactionID:     fmt.Sprintf("tx%d", i+2),
			Amount:            50.0,
		}

		games = append(games, game)
	}

	// Mock select game results for the user
	mockQuerier.On("SelectGameResultsByUser", ctx, userID, entity.ValidationStatusPending).Return(games, nil)

	// Mock update game result
	mockQuerier.On("UpdateGameResult", ctx, mock.Anything, mock.Anything, entity.ValidationStatusAccepted).Times((totalEntries + 1) - totalGamesToCancel).Return(nil)
	mockQuerier.On("UpdateGameResult", ctx, mock.Anything, mock.Anything, entity.ValidationStatusCanceled).Times(totalGamesToCancel).Return(nil)

	// Mock update user balance
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userID, balanceExpectedAdjustment, true).Return(nil)

	// Mock transaction operations
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error")).Return(nil)

	err := instance.ValidateGameResults(ctx, totalGamesToCancel)

	assert.NoError(t, err, "ValidateGameResults should not return an error")
	mockQuerier.AssertExpectations(t)
}

func TestValidateGameResultsSelectUsersError(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()

	// Mock select users with pending validation error
	mockQuerier.On("SelectUsersByValidationStatus", ctx, false).Return(nil, errors.New("database error"))

	err := instance.ValidateGameResults(ctx, 1)

	assert.Error(t, err, "ValidateGameResults should return an error on SelectUsersByValidationStatus")
	mockQuerier.AssertExpectations(t)
}

func TestValidateGameResultsSelectGameResultsError(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()

	instance := NewGameResultDAO(mockQuerier)

	ctx := context.TODO()

	userID := uuid.New()

	// Mock select users with pending validation
	mockQuerier.On("SelectUsersByValidationStatus", ctx, false).Return([]entity.User{
		{
			ID:      userID,
			Balance: 100.0,
		},
	}, nil)

	// Mock select game results for the user error
	mockQuerier.On("SelectGameResultsByUser", ctx, userID, entity.ValidationStatusPending).Return(nil, errors.New("database error"))

	err := instance.ValidateGameResults(ctx, 1)

	assert.Error(t, err, "ValidateGameResults should return an error on SelectGameResultsByUser")
	mockQuerier.AssertExpectations(t)
}
