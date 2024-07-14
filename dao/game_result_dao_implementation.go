package dao

import (
	"context"
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

func NewGameResultDAO(querier database.Querier) *gameResultDAO {
	dm := gameResultDAO{
		querier: querier,
	}
	return &dm
}

// CreateGameResult creates a new game result in the database.
func (dm *gameResultDAO) CreateGameResult(ctx context.Context, userId uuid.UUID, gameStatus entity.GameStatus, amount float64, transactionSource entity.TransactionSource, transactionID string) (*entity.GameResult, error) {

	// Validate the existence of the transaction ID.
	exists, err := dm.querier.CheckTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, entity.ErrTransactionIdExists
	}

	// Validate the user.
	user, err := dm.querier.SelectUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	// Check if the user has a negative balance.
	if gameStatus == entity.GameStatusLost && user.Balance < amount {
		return nil, entity.ErrUserNegativeBalance
	}

	// Define the new balance.
	balance := user.Balance
	if gameStatus == entity.GameStatusWin {
		balance = user.Balance + amount
	} else {
		balance = user.Balance - amount
	}

	// Define the game result.
	gameResult := entity.GameResult{
		UserID:            userId,
		GameStatus:        gameStatus,
		ValidationStatus:  entity.ValidationStatusPending,
		TransactionSource: transactionSource,
		TransactionID:     transactionID,
		Amount:            amount,
		CreatedAt:         time.Now(),
	}

	// Starts a new transaction
	err = dm.querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {

		// Lock the user row
		err := dm.querier.LockUserRow(ctx, *txn, userId)
		if err != nil {
			return err
		}

		// Create the game result
		id, err := dm.querier.InsertGameResult(ctx, *txn, gameResult)
		if err != nil {
			return err
		}
		gameResult.ID = id

		// Update the user balance
		err = dm.querier.UpdateUserBalance(ctx, *txn, userId, balance, false)
		if err != nil {
			return err
		}

		// Commit the transaction
		// Success, continue with the transaction commit
		return nil
	})

	// In case of error in the transactional operation, log and exit
	if err != nil {
		log.Printf("error performing game result db transaction: %v", err)
		return nil, entity.ErrCreatingGameResult
	}

	return &gameResult, nil
}

// ValidateGameResults validates the game results in post-game processing.
func (dm *gameResultDAO) ValidateGameResults(ctx context.Context, totalGamesToCancel int) error {

	// Retrieve users pending validation
	users, err := dm.querier.SelectUsersByValidationStatus(ctx, false)
	if err != nil {
		return err
	}
	log.Printf("validating latest users. Total users %d", len(users))

	// Iterate over the pending users
	for _, user := range users {
		// Retrieve the user's game results
		gameResults, err := dm.querier.SelectGameResultsByUser(ctx, user.ID, entity.ValidationStatusPending)
		if err != nil {
			return err
		}

		// Reset counter
		totalTransactionsCanceled := 0

		// Initial balance
		balance := user.Balance

		// Starts a new transaction
		err = dm.querier.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			// Lock the user row
			err := dm.querier.LockUserRow(ctx, *txn, user.ID)
			if err != nil {
				return err
			}

			// Iterate over the user's pending game results
			for _, gameResult := range gameResults {

				// Check if the total has been achieved and the entry ID is odd
				if totalTransactionsCanceled < totalGamesToCancel &&
					gameResult.ShouldBeCanceled() {

					// Cancel the game result
					err := dm.querier.UpdateGameResult(ctx, *txn, gameResult.ID, entity.ValidationStatusCanceled)
					if err != nil {
						return err
					}

					// Calculate the balance reversal
					if gameResult.GameStatus == entity.GameStatusWin {
						balance -= gameResult.Amount
					} else {
						balance += gameResult.Amount
					}

					// Increment the counter
					totalTransactionsCanceled++

				} else {
					// Approve everything else
					err := dm.querier.UpdateGameResult(ctx, *txn, gameResult.ID, entity.ValidationStatusAccepted)
					if err != nil {
						return err
					}
				}
			}

			// Update the user balance
			// Also set user balance as validated
			err = dm.querier.UpdateUserBalance(ctx, *txn, user.ID, balance, true)
			if err != nil {
				return err
			}

			// Commit the transaction
			// Success, continue with the transaction commit
			return nil
		})

		// In case of error in the transactional operation, log and exit
		if err != nil {
			return err
		}

		log.Printf("%d game results cancelled for user %s.", totalTransactionsCanceled, user.ID)
	}

	return nil
}
