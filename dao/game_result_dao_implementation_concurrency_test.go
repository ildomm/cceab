package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/database"
	"github.com/ildomm/cceab/entity"
	"github.com/ildomm/cceab/test_helpers"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestCreateGameResultSequenceOnMock(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	ctx := context.Background()

	instance := NewGameResultDAO(mockQuerier)

	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame
	transactionID := "unique-transaction-id"

	// Mock successful interactions
	mockQuerier.On("CheckTransactionID", ctx, mock.Anything).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return() // no fake results
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error"))
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userId).Return(nil)
	mockQuerier.On("InsertGameResult", ctx, mock.Anything, mock.Anything).Return(uuid.New(), nil)
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userId, mock.Anything, false).Return(nil)

	// Give to the mock a user with a balance of 0
	mockQuerier.WithTransaction(ctx, func(txn *sqlx.Tx) error {
		mockQuerier.UpdateUserBalance(ctx, *txn, userId, 0, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [1000]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)
	expectedBalance := amount * float64(totalInjected)

	for range toInjectTotalEntries {
		transactionID = uuid.New().String()
		_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, transactionID)
		assert.NoError(t, err)
	}

	// Basic mockers expectations check
	mockQuerier.AssertExpectations(t)

	// Count the game results
	assert.Equal(t, mockQuerier.GameCount(), totalInjected)

	// Compare the use balance
	user, err := mockQuerier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}

func TestCreateGameResultConcurrentOnMock(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	ctx := context.Background()

	instance := NewGameResultDAO(mockQuerier)

	userId := uuid.New()
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame

	// Mock successful interactions
	mockQuerier.On("CheckTransactionID", ctx, mock.Anything).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return() // no fake results
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error"))
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userId).Return(nil)
	mockQuerier.On("InsertGameResult", ctx, mock.Anything, mock.Anything).Return(uuid.New(), nil)
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userId, mock.Anything, false).Return(nil)

	// Give to the mock a user with a balance of 0
	mockQuerier.WithTransaction(ctx, func(txn *sqlx.Tx) error {
		mockQuerier.UpdateUserBalance(ctx, *txn, userId, 0, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [1000]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)
	expectedBalance := amount * float64(totalInjected)

	wg := sync.WaitGroup{}
	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()
			_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	// Wait for all workers to complete processing
	wg.Wait()

	// Basic mockers expectations check
	mockQuerier.AssertExpectations(t)

	// Count the game results
	assert.Equal(t, mockQuerier.GameCount(), totalInjected)

	// Compare the use balance
	user, err := mockQuerier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}

func TestCreateGameResultConcurrentWinLostOnMock(t *testing.T) {
	mockQuerier := test_helpers.NewMockQuerier()
	ctx := context.Background()

	instance := NewGameResultDAO(mockQuerier)

	userId := uuid.New()
	amount := 100.0
	transactionSource := entity.TransactionSourceGame

	// Mock successful interactions
	mockQuerier.On("CheckTransactionID", ctx, mock.Anything).Return(false, nil)
	mockQuerier.On("SelectUser", ctx, userId).Return() // no fake results
	mockQuerier.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(*sqlx.Tx) error"))
	mockQuerier.On("LockUserRow", ctx, mock.Anything, userId).Return(nil)
	mockQuerier.On("InsertGameResult", ctx, mock.Anything, mock.Anything).Return(uuid.New(), nil)
	mockQuerier.On("UpdateUserBalance", ctx, mock.Anything, userId, mock.Anything, false).Return(nil)

	// Give to the mock a user with a balance of 1000
	mockQuerier.WithTransaction(ctx, func(txn *sqlx.Tx) error {

		// Must start with some balance, unless the user will have a negative balance for the first
		// entity.GameStatusLost hit
		mockQuerier.UpdateUserBalance(ctx, *txn, userId, 1000, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [100]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)

	// To equalize the balance
	expectedBalance := float64(1000)

	wg := sync.WaitGroup{}
	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()

			_, err := instance.CreateGameResult(ctx, userId, entity.GameStatusWin, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()
			_, err := instance.CreateGameResult(ctx, userId, entity.GameStatusLost, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	// Wait for all workers to complete processing
	wg.Wait()

	// Basic mockers expectations check
	mockQuerier.AssertExpectations(t)

	// Count the game results
	assert.Equal(t, mockQuerier.GameCount(), totalInjected*2)

	// Compare the use balance
	user, err := mockQuerier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}

func setupTestQuerier(t *testing.T) (context.Context, func(t *testing.T), *database.PostgresQuerier) {
	testDB := test_helpers.NewTestDatabase(t)
	ctx := context.Background()
	q, err := database.NewPostgresQuerier(ctx, testDB.ConnectionString(t)+"?sslmode=disable")
	require.NoError(t, err)

	return ctx, func(t *testing.T) {
		testDB.Close(t)
	}, q
}

func TestCreateGameResultSequenceOnDB(t *testing.T) {
	ctx, teardownTest, querier := setupTestQuerier(t)
	defer teardownTest(t)

	instance := NewGameResultDAO(querier)

	userId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame

	// Give to the mock a user with a balance of 0
	querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {
		querier.UpdateUserBalance(ctx, *txn, userId, 0, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [100]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)
	expectedBalance := amount * float64(totalInjected)

	for range toInjectTotalEntries {
		_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, uuid.New().String())
		assert.NoError(t, err)
	}

	// Count the game results
	games, err := querier.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusPending)
	assert.NoError(t, err)
	assert.Len(t, games, totalInjected)

	// Compare the use balance
	user, err := querier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}

func TestCreateGameResultConcurrentOnDB(t *testing.T) {
	ctx, teardownTest, querier := setupTestQuerier(t)
	defer teardownTest(t)

	instance := NewGameResultDAO(querier)

	userId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	gameStatus := entity.GameStatusWin
	amount := 100.0
	transactionSource := entity.TransactionSourceGame

	// Give to the mock a user with a balance of 0
	querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {
		querier.UpdateUserBalance(ctx, *txn, userId, 0, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [100]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)
	expectedBalance := amount * float64(totalInjected)

	wg := sync.WaitGroup{}
	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()
			_, err := instance.CreateGameResult(ctx, userId, gameStatus, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	// Wait for all workers to complete processing
	wg.Wait()

	// Count the game results
	games, err := querier.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusPending)
	assert.NoError(t, err)
	assert.Len(t, games, totalInjected)

	// Compare the use balance
	user, err := querier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}

func TestCreateGameResultConcurrentWinLostOnDB(t *testing.T) {
	ctx, teardownTest, querier := setupTestQuerier(t)
	defer teardownTest(t)

	instance := NewGameResultDAO(querier)

	userId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	amount := 100.0
	transactionSource := entity.TransactionSourceGame

	// Give to the mock a user with a balance of 1000
	querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {

		// Must start with some balance, unless the user will have a negative balance for the first
		// entity.GameStatusLost hit
		querier.UpdateUserBalance(ctx, *txn, userId, 1000, false)
		return nil
	})

	// Inject many game results
	toInjectTotalEntries := [100]int{} //nolint:all
	totalInjected := len(toInjectTotalEntries)

	// To equalize the balance
	expectedBalance := float64(1000)

	wg := sync.WaitGroup{}
	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()

			_, err := instance.CreateGameResult(ctx, userId, entity.GameStatusWin, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	for range toInjectTotalEntries {
		wg.Add(1)

		// A go routine for each game result
		go func() {
			defer wg.Done()
			_, err := instance.CreateGameResult(ctx, userId, entity.GameStatusLost, amount, transactionSource, uuid.New().String())
			assert.NoError(t, err)
		}()
	}

	// Wait for all workers to complete processing
	wg.Wait()

	// Count the game results
	games, err := querier.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusPending)
	assert.NoError(t, err)
	assert.Len(t, games, totalInjected*2)

	// Compare the use balance
	user, err := querier.SelectUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, user.Balance)
}
