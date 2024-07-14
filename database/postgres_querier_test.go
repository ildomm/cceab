package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/ildomm/cceab/test_helpers"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestPostgresQuerier(t *testing.T) {
	testDB := test_helpers.NewTestDatabase(t)
	dbURL := testDB.ConnectionString(t) + "?sslmode=disable"

	ctx := context.Background()

	t.Run("NewPostgresQuerier_Success", func(t *testing.T) {
		querier, err := NewPostgresQuerier(ctx, dbURL)
		require.NoError(t, err)
		require.NotNil(t, querier)

		defer querier.Close()

		assert.NotNil(t, querier.dbConn)

		// Check if at least one migration has run by querying the database
		var extensionExists bool
		err = querier.dbConn.Get(&extensionExists, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp')")
		require.NoError(t, err)
		assert.True(t, extensionExists, "uuid-ossp extension should exist")

		// Check the number of migration files in the folder
		migrationFiles, err := fs.ReadDir("migrations")
		require.NoError(t, err)

		// Count the number of migration files
		expectedNumMigrations := len(migrationFiles) / 2 // Each migration file has a corresponding .down.sql file

		// Query the "schema_migrations_service_file_loader" table to get the version
		var version string
		err = querier.dbConn.Get(&version, "SELECT version FROM schema_migrations")
		require.NoError(t, err)

		// Convert the version to an integer
		versionInt, err := strconv.Atoi(strings.TrimSpace(version))
		require.NoError(t, err)

		// Compare the number of migrations with the version in the database
		assert.Equal(t, expectedNumMigrations, versionInt, fmt.Sprintf("Number of migrations should match the version in the database. Expected: %d, Actual: %d", expectedNumMigrations, versionInt))
	})

	t.Run("NewPostgresQuerier_InvalidURL", func(t *testing.T) {
		_, err := NewPostgresQuerier(ctx, "invalid-url")
		require.Error(t, err)
	})
}

func setupTestQuerier(t *testing.T) (context.Context, func(t *testing.T), *PostgresQuerier) {
	testDB := test_helpers.NewTestDatabase(t)
	ctx := context.Background()
	q, err := NewPostgresQuerier(ctx, testDB.ConnectionString(t)+"?sslmode=disable")
	require.NoError(t, err)

	return ctx, func(t *testing.T) {
		testDB.Close(t)
	}, q
}

func TestDatabaseWithTransaction(t *testing.T) {
	ctx, teardownTest, q := setupTestQuerier(t)
	defer teardownTest(t)

	userId, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)

	t.Run("LockUserRow_Success", func(t *testing.T) {
		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			err := q.LockUserRow(ctx, *txn, userId)
			require.NoError(t, err)

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("InsertGameResult_Success", func(t *testing.T) {
		gameResult := entity.GameResult{
			UserID:            userId,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusPending,
			TransactionSource: entity.TransactionSourceServer,
			TransactionID:     "anything",
			Amount:            10,
			CreatedAt:         time.Now(),
		}

		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			id, err := q.InsertGameResult(ctx, *txn, gameResult)
			require.NoError(t, err)

			gameResult.ID = id

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("UpdateUserBalance_Success", func(t *testing.T) {

		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			err := q.UpdateUserBalance(ctx, *txn, userId, 100, true)
			require.NoError(t, err)

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)

		// Check the new balance
		user, err := q.SelectUser(ctx, userId)
		require.NoError(t, err)
		require.Equal(t, 100.0, user.Balance)
	})
}

func TestDatabaseBasicOperations(t *testing.T) {
	ctx, teardownTest, q := setupTestQuerier(t)
	defer teardownTest(t)

	userId, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)

	t.Run("SelectUser_Success", func(t *testing.T) {
		user, err := q.SelectUser(ctx, userId)
		require.NoError(t, err)
		require.Equal(t, userId, user.ID)
	})

	t.Run("SelectUsersByValidationStatus_Success", func(t *testing.T) {
		users, err := q.SelectUsersByValidationStatus(ctx, false)
		require.NoError(t, err)
		require.Len(t, users, 4)
	})

	t.Run("CheckTransactionID_None", func(t *testing.T) {
		exist, err := q.CheckTransactionID(ctx, "anything")
		require.NoError(t, err)
		require.False(t, exist)
	})

	t.Run("CheckTransactionID_One", func(t *testing.T) {
		gameResult := entity.GameResult{
			UserID:            userId,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusPending,
			TransactionSource: entity.TransactionSourceServer,
			TransactionID:     "anything",
			Amount:            10,
			CreatedAt:         time.Now(),
		}

		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			id, err := q.InsertGameResult(ctx, *txn, gameResult)
			require.NoError(t, err)

			gameResult.ID = id

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)

		exist, err := q.CheckTransactionID(ctx, "anything")
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("SelectGameResultsByUser_Success", func(t *testing.T) {
		gameResult := entity.GameResult{
			UserID:            userId,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusCanceled,
			TransactionSource: entity.TransactionSourceServer,
			TransactionID:     "anything",
			Amount:            10,
			CreatedAt:         time.Now(),
		}

		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			id, err := q.InsertGameResult(ctx, *txn, gameResult)
			require.NoError(t, err)

			gameResult.ID = id

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)

		gameResults, err := q.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusCanceled)
		require.NoError(t, err)
		require.Len(t, gameResults, 1)

		gameResults, err = q.SelectGameResultsByUser(ctx, userId, entity.ValidationStatusAccepted)
		require.NoError(t, err)
		require.Len(t, gameResults, 0)
	})

	t.Run("UpdateGameResult_Success", func(t *testing.T) {
		gameResult := entity.GameResult{
			UserID:            userId,
			GameStatus:        entity.GameStatusWin,
			ValidationStatus:  entity.ValidationStatusPending,
			TransactionSource: entity.TransactionSourceServer,
			TransactionID:     "anything",
			Amount:            10,
			CreatedAt:         time.Now(),
		}

		// Start a transaction that is expected to WORK
		err := q.WithTransaction(ctx, func(txn *sqlx.Tx) error {

			id, err := q.InsertGameResult(ctx, *txn, gameResult)
			require.NoError(t, err)
			gameResult.ID = id

			err = q.UpdateGameResult(ctx, *txn, gameResult.ID, entity.ValidationStatusAccepted)
			require.NoError(t, err)

			// No error, then the db commit() will happen
			return nil
		})
		require.NoError(t, err)

	})
}
