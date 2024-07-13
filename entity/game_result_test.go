package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGameStatus_Scan(t *testing.T) {
	var status GameStatus
	err := status.Scan("win")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if status != GameStatusWin {
		t.Errorf("Expected %v, got %v", GameStatusWin, status)
	}
}

func TestGameStatus_Value(t *testing.T) {
	status := GameStatusWin
	val, err := status.Value()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != "win" {
		t.Errorf("Expected 'win', got %v", val)
	}
}

func TestValidationStatus_Scan(t *testing.T) {
	var status ValidationStatus
	err := status.Scan("pending")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if status != ValidationStatusPending {
		t.Errorf("Expected %v, got %v", ValidationStatusPending, status)
	}
}

func TestValidationStatus_Value(t *testing.T) {
	status := ValidationStatusPending
	val, err := status.Value()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != "pending" {
		t.Errorf("Expected 'pending', got %v", val)
	}
}

func TestTransactionSource_Scan(t *testing.T) {
	var source TransactionSource
	err := source.Scan("game")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if source != TransactionSourceGame {
		t.Errorf("Expected %v, got %v", TransactionSourceGame, source)
	}
}

func TestTransactionSource_Value(t *testing.T) {
	source := TransactionSourceGame
	val, err := source.Value()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != "game" {
		t.Errorf("Expected 'game', got %v", val)
	}
}

func TestGameResult(t *testing.T) {
	id := 1
	userID := uuid.New()
	gameStatus := GameStatusWin
	validationStatus := ValidationStatusPending
	transactionSource := TransactionSourceGame
	transactionID := "tx123"
	amount := 10.15
	createdAt := time.Now()

	gameResult := GameResult{
		ID:                id,
		UserID:            userID,
		GameStatus:        gameStatus,
		ValidationStatus:  validationStatus,
		TransactionSource: transactionSource,
		TransactionID:     transactionID,
		Amount:            amount,
		CreatedAt:         createdAt,
	}

	if gameResult.ID != id {
		t.Errorf("Expected ID %v, got %v", id, gameResult.ID)
	}
	if gameResult.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, gameResult.UserID)
	}
	if gameResult.GameStatus != gameStatus {
		t.Errorf("Expected GameStatus %v, got %v", gameStatus, gameResult.GameStatus)
	}
	if gameResult.ValidationStatus != validationStatus {
		t.Errorf("Expected ValidationStatus %v, got %v", validationStatus, gameResult.ValidationStatus)
	}
	if gameResult.TransactionSource != transactionSource {
		t.Errorf("Expected TransactionSource %v, got %v", transactionSource, gameResult.TransactionSource)
	}
	if gameResult.TransactionID != transactionID {
		t.Errorf("Expected TransactionID %v, got %v", transactionID, gameResult.TransactionID)
	}
	if gameResult.Amount != amount {
		t.Errorf("Expected Amount %v, got %v", amount, gameResult.Amount)
	}
	if !gameResult.CreatedAt.Equal(createdAt) {
		t.Errorf("Expected CreatedAt %v, got %v", createdAt, gameResult.CreatedAt)
	}
}
