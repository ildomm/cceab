package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"time"
)

type PostgresQuerier struct {
	dbURL  string
	dbConn *sqlx.DB
	ctx    context.Context
}

func NewPostgresQuerier(ctx context.Context, url string) (*PostgresQuerier, error) {
	querier := PostgresQuerier{dbURL: url, ctx: ctx}

	_, err := pgx.ParseConfig(url)
	if err != nil {
		return &querier, err
	}

	// Open the connection using the DataDog wrapper methods
	querier.dbConn, err = sqlx.Open("pgx", url)
	if err != nil {
		return &querier, err
	}
	log.Print("opened database connection")

	// Ping the database to check that the connection is actually working
	err = querier.dbConn.Ping()
	if err != nil {
		return &querier, err
	}

	// Migrate the database
	err = querier.migrate()
	if err != nil {
		return &querier, err
	}
	log.Print("database migration complete")

	return &querier, nil
}

func (q *PostgresQuerier) Close() {
	q.dbConn.Close()
	log.Print("closed database connection")
}

var (
	//go:embed migrations/*.sql
	fs           embed.FS
	ErrorNilUUID = errors.New("UUID is nil")
)

func (q *PostgresQuerier) migrate() error {

	// Amend the database URl with custom parameter so that we can specify the
	// table name to be used to hold database migration state
	migrationsURL, err := q.migrationsURL()
	if err != nil {
		return err
	}

	// Load the migrations from our embedded resources
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	// Use a custom table name for schema migrations
	m, err := migrate.NewWithSourceInstance("iofs", d, migrationsURL)
	if err != nil {
		return err
	}

	// Migrate all the way up ...
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

const (
	CustomMigrationParam = "x-migrations-table"
	CustomMigrationValue = "schema_migrations"
)

func (q *PostgresQuerier) migrationsURL() (string, error) {
	url, err := url.Parse(q.dbURL)
	if err != nil {
		return "", err
	}

	// Add the new Query parameter that specifies the table name for the migrations
	values := url.Query()
	values.Add(CustomMigrationParam, CustomMigrationValue)

	// Replace the Query parameters in the original URL & return
	url.RawQuery = values.Encode()
	return url.String(), nil
}

////////////////////////////////// Database Querier standard operations /////////////////////////////////////////////////////////

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func (q *PostgresQuerier) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) (err error) {

	// Starting database transaction
	tx, err := q.dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			err = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	// The database transaction starts taking effect here
	err = fn(tx)
	return err
}

////////////////////////////////// Database Querier domain operations /////////////////////////////////////////////////////////

const insertGameResultSQL = `
	INSERT INTO game_results ( user_id, game_status, validation_status, transaction_source, transaction_id, amount, created_at)
	VALUES                 ( $1,      $2,          $3,                $4,                 $5,             $6,     $7)
	RETURNING id`

func (q *PostgresQuerier) InsertGameResult(ctx context.Context, txn sqlx.Tx, gameResult entity.GameResult) (int, error) {
	var id int

	err := txn.GetContext(
		ctx,
		&id,
		insertGameResultSQL,
		gameResult.UserID,
		gameResult.GameStatus,
		gameResult.ValidationStatus,
		gameResult.TransactionSource,
		gameResult.TransactionID,
		gameResult.Amount,
		gameResult.CreatedAt)

	return id, err
}

const lockUserRowStep1SQL = `LOCK TABLE users IN ROW EXCLUSIVE MODE;`
const lockUserRowStep2SQL = `SELECT * FROM users WHERE id = $1 FOR UPDATE;`

func (q *PostgresQuerier) LockUserRow(ctx context.Context, txn sqlx.Tx, userId uuid.UUID) error {
	_, err := txn.ExecContext(ctx, lockUserRowStep1SQL)
	if err != nil {
		return err
	}

	_, err = txn.ExecContext(ctx, lockUserRowStep2SQL, userId)
	if err != nil {
		return err
	}

	return err
}

const selectUserSQL = `SELECT * FROM users WHERE id = $1`

func (q *PostgresQuerier) SelectUser(ctx context.Context, userId uuid.UUID) (*entity.User, error) {
	var user entity.User

	err := q.dbConn.GetContext(
		ctx,
		&user,
		selectUserSQL,
		userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

const selectUserByValidationStatusSQL = `SELECT * FROM users WHERE games_result_validated = $1`

func (q *PostgresQuerier) SelectUsersByValidationStatus(ctx context.Context, validationStatus bool) ([]entity.User, error) {
	var users []entity.User

	err := q.dbConn.SelectContext(
		ctx,
		&users,
		selectUserByValidationStatusSQL,
		validationStatus)

	return users, err
}

const selectCheckTransactionSQL = `SELECT count(*) FROM game_results WHERE transaction_id = $1`

func (q *PostgresQuerier) CheckTransactionID(ctx context.Context, transactionId string) (bool, error) {

	row := q.dbConn.QueryRowContext(ctx, selectCheckTransactionSQL, transactionId)
	var count int64

	err := row.Scan(&count)

	if count > 0 {
		return true, err
	} else {
		return false, err
	}
}

const selectGameResultsByUserSQL = `SELECT * FROM game_results WHERE user_id = $1 AND validation_status = $2 ORDER BY created_at DESC`

func (q *PostgresQuerier) SelectGameResultsByUser(ctx context.Context, userId uuid.UUID, validationStatus entity.ValidationStatus) ([]entity.GameResult, error) {
	var gameResults []entity.GameResult

	err := q.dbConn.SelectContext(
		ctx,
		&gameResults,
		selectGameResultsByUserSQL,
		userId,
		validationStatus)

	return gameResults, err
}

const updateUserSQL = `
	UPDATE users
	SET 
		balance = :balance,
		games_result_validated = :games_result_validated,
		last_game_result_at = :last_game_result_at
	WHERE id = :id`

func (q *PostgresQuerier) UpdateUserBalance(ctx context.Context, txn sqlx.Tx, userId uuid.UUID, balance float64, validationStatus bool) error {
	user := entity.User{
		ID:                   userId,
		Balance:              balance,
		GamesResultValidated: sql.NullBool{Bool: validationStatus, Valid: true},
		LastGameResultAt:     sql.NullTime{Time: time.Now(), Valid: true},
	}

	_, err := txn.NamedExecContext(ctx, updateUserSQL, user)

	return err
}

const updateGameResultSQL = `
	UPDATE game_results
	SET 
		validation_status = :validation_status
	WHERE id = :id`

func (q *PostgresQuerier) UpdateGameResult(ctx context.Context, txn sqlx.Tx, gameResultId int, validationStatus entity.ValidationStatus) error {
	gameResult := entity.GameResult{
		ID:               gameResultId,
		ValidationStatus: validationStatus,
	}

	_, err := txn.NamedExecContext(ctx, updateGameResultSQL, gameResult)

	return err
}
