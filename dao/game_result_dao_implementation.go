package dao

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/database"
	"github.com/ildomm/cceab/entity"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type gameResultDAO struct {
	querier database.Querier
}

// NewGameResultDAO creates a new game result DAO
func NewGameResultDAO(querier database.Querier) *gameResultDAO {
	return &gameResultDAO{querier: querier}
}

// CreateGameResult creates a new game result
// It validates the transaction and updates the user balance
// It returns the created game result
// It returns an error if the transaction is invalid or if there is an error creating the game result
func (dm *gameResultDAO) CreateGameResult(ctx context.Context, userId uuid.UUID, gameStatus entity.GameStatus, amount float64, transactionSource entity.TransactionSource, transactionID string) (*entity.GameResult, error) {

	// Check the transaction and its related user
	balance := 0.0
	if user, err := dm.validateTransaction(ctx, userId, gameStatus, amount, transactionID); err != nil {
		return nil, err
	} else {
		balance = dm.calculateNewBalance(user.Balance, gameStatus, amount)
	}

	gameResult := entity.GameResult{
		UserID:            userId,
		GameStatus:        gameStatus,
		ValidationStatus:  entity.ValidationStatusPending,
		TransactionSource: transactionSource,
		TransactionID:     transactionID,
		Amount:            amount,
		CreatedAt:         time.Now(),
	}

	// Perform the whole operation inside a db transaction
	err := dm.querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {
		if err := dm.persistGameResultTransaction(ctx, txn, userId, &gameResult, balance); err != nil {
			log.Printf("error persisting game result: %v", err)
			return err
		}

		// Commit the transaction
		// Success, continue with the transaction commit
		return nil
	})
	if err != nil {
		log.Printf("error performing game result db transaction: %v", err)
		return nil, entity.ErrCreatingGameResult
	}

	return &gameResult, nil
}

// validateTransaction validates the transaction
// It returns the user if the transaction is valid
func (dm *gameResultDAO) validateTransaction(ctx context.Context, userId uuid.UUID, gameStatus entity.GameStatus, amount float64, transactionID string) (*entity.User, error) {
	exists, err := dm.querier.CheckTransactionID(ctx, transactionID)
	if err != nil {
		log.Printf("error locating transaction: %v", err)
		return nil, err
	}
	if exists {
		return nil, entity.ErrTransactionIdExists
	}

	user, err := dm.querier.SelectUser(ctx, userId)
	if err != nil {
		log.Printf("error locating user: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	// No negative balance allowed
	if gameStatus == entity.GameStatusLost && user.Balance < amount {
		return nil, entity.ErrUserNegativeBalance
	}

	return user, nil
}

// calculateNewBalance calculates the new balance based on the game status
func (dm *gameResultDAO) calculateNewBalance(currentBalance float64, gameStatus entity.GameStatus, amount float64) float64 {
	if gameStatus == entity.GameStatusWin {
		return currentBalance + amount
	}
	return currentBalance - amount
}

// persistGameResultTransaction persists the game result transaction
func (dm *gameResultDAO) persistGameResultTransaction(ctx context.Context, txn *sqlx.Tx, userId uuid.UUID, gameResult *entity.GameResult, balance float64) error {

	// No other processes can update the user until end of this transaction
	if err := dm.querier.LockUserRow(ctx, *txn, userId); err != nil {
		return fmt.Errorf("locking user row: %w", err)
	}

	id, err := dm.querier.InsertGameResult(ctx, *txn, *gameResult)
	if err != nil {
		return fmt.Errorf("inserting game result: %w", err)
	}
	gameResult.ID = id

	if err := dm.querier.UpdateUserBalance(ctx, *txn, userId, balance, false); err != nil {
		return fmt.Errorf("updating user balance: %w", err)
	}

	return nil
}

// ValidateGameResults validates the game results
// It cancels the game results that should be canceled and approves the rest
// It returns an error if there is an error validating the game results
func (dm *gameResultDAO) ValidateGameResults(ctx context.Context, totalGamesToCancel int) error {

	// Select the latest users that have not been validated
	users, err := dm.querier.SelectUsersByValidationStatus(ctx, false)
	if err != nil {
		return fmt.Errorf("selecting users by validation status: %w", err)
	}
	log.Printf("validating latest users. Total users %d", len(users))

	for _, user := range users {
		if err := dm.validateUserGameResults(ctx, user, totalGamesToCancel); err != nil {
			return fmt.Errorf("validating user game results for user %s: %w", user.ID, err)
		}
	}
	return nil
}

// validateUserGameResults validates the game results for a user
func (dm *gameResultDAO) validateUserGameResults(ctx context.Context, user entity.User, totalGamesToCancel int) error {

	// Select the game results for the user that are pending validation
	gameResults, err := dm.querier.SelectGameResultsByUser(ctx, user.ID, entity.ValidationStatusPending)
	if err != nil {
		return fmt.Errorf("selecting game results by user: %w", err)
	}

	balance := user.Balance
	totalTransactionsCanceled := 0

	// Perform the whole operation inside a db transaction
	err = dm.querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {

		// No other processes can update the user until end of this transaction
		if err := dm.querier.LockUserRow(ctx, *txn, user.ID); err != nil {
			return fmt.Errorf("locking user row: %w", err)
		}

		// Check all the game results, until:
		// - All the transactions to cancel have been canceled, based on the limit (totalGamesToCancel)
		// - All the rest of the transactions have been approved
		for _, gameResult := range gameResults {
			if totalTransactionsCanceled < totalGamesToCancel && gameResult.ShouldBeCanceled() {
				if err := dm.cancelGameResult(ctx, txn, gameResult, &balance); err != nil {
					return fmt.Errorf("canceling game result: %w", err)
				}
				totalTransactionsCanceled++
			} else {
				if err := dm.approveGameResult(ctx, txn, gameResult); err != nil {
					return fmt.Errorf("approving game result: %w", err)
				}
			}
		}

		// Reset the user balance
		if err := dm.querier.UpdateUserBalance(ctx, *txn, user.ID, balance, true); err != nil {
			return fmt.Errorf("updating user balance: %w", err)
		}

		// Commit the transaction
		// Success, continue with the transaction commit
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("%d game results cancelled for user %s.", totalTransactionsCanceled, user.ID)
	return nil
}

// cancelGameResult cancels the game result
// Calculates the new balance based on the game status
func (dm *gameResultDAO) cancelGameResult(ctx context.Context, txn *sqlx.Tx, gameResult entity.GameResult, balance *float64) error {
	if err := dm.querier.UpdateGameResult(ctx, *txn, gameResult.ID, entity.ValidationStatusCanceled); err != nil {
		return fmt.Errorf("updating game result to canceled: %w", err)
	}

	if gameResult.GameStatus == entity.GameStatusWin {
		*balance -= gameResult.Amount
	} else {
		*balance += gameResult.Amount
	}

	return nil
}

// approveGameResult approves the game result
func (dm *gameResultDAO) approveGameResult(ctx context.Context, txn *sqlx.Tx, gameResult entity.GameResult) error {
	return dm.querier.UpdateGameResult(ctx, *txn, gameResult.ID, entity.ValidationStatusAccepted)
}
