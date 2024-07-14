package test_helpers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"sync"
)

// MockQuerier is a mock type for the Querier type
type MockQuerier struct {
	mock.Mock
	lock sync.Mutex

	keys map[string]map[string]interface{}
}

// NewMockQuerier creates a new instance of MockQuerier
func NewMockQuerier() *MockQuerier {
	mocked := &MockQuerier{
		keys: make(map[string]map[string]interface{}),
	}

	mocked.keys["game_results"] = make(map[string]interface{})

	return mocked
}

func (m *MockQuerier) Close() {
	m.Called()
}

func (m *MockQuerier) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) (err error) {
	m.Called(ctx, fn)

	txn := new(sqlx.Tx)
	err = fn(txn)

	return err
}

func (m *MockQuerier) InsertGameResult(ctx context.Context, txn sqlx.Tx, gameResult entity.GameResult) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, txn, gameResult)
	if arg := args.Get(0); arg != nil {

		id := rand.Intn(1000)
		gameResult.ID = id

		m.keys["game_results"][fmt.Sprint(id)] = gameResult

		return id, nil
	}
	return 0, args.Error(1)
}

func (m *MockQuerier) LockUserRow(ctx context.Context, txn sqlx.Tx, userId uuid.UUID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, txn, userId)
	return args.Error(0)
}

func (m *MockQuerier) SelectUser(ctx context.Context, userId uuid.UUID) (*entity.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, userId)
	if arg := args.Get(0); arg != nil {
		return arg.(*entity.User), nil
	}
	return nil, args.Error(1)
}

func (m *MockQuerier) SelectUsersByValidationStatus(ctx context.Context, validationStatus bool) ([]entity.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, validationStatus)
	if arg := args.Get(0); arg != nil {
		return arg.([]entity.User), nil
	}
	return nil, args.Error(1)
}

func (m *MockQuerier) CheckTransactionID(ctx context.Context, transactionId string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, transactionId)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) GameResultExists(ctx context.Context, transactionId string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, transactionId)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) SelectGameResultsByUser(ctx context.Context, userId uuid.UUID, validationStatus entity.ValidationStatus) ([]entity.GameResult, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, userId, validationStatus)
	if len(args) > 0 {
		if args[1] != nil {
			return nil, args.Error(1)
		} else {
			return args.Get(0).([]entity.GameResult), nil
		}
	}

	var games []entity.GameResult
	for _, gameResult := range m.keys["game_results"] {
		games = append(games, gameResult.(entity.GameResult))
	}
	return games, nil

}

func (m *MockQuerier) UpdateUserBalance(ctx context.Context, txn sqlx.Tx, userId uuid.UUID, amount float64, validationStatus bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, txn, userId, amount, validationStatus)
	return args.Error(0)
}

func (m *MockQuerier) UpdateGameResult(ctx context.Context, txn sqlx.Tx, gameResultId int, validationStatus entity.ValidationStatus) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(ctx, txn, gameResultId, validationStatus)
	return args.Error(0)
}
