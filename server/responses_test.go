package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestWriteInternalError tests the WriteInternalError function.
func TestWriteInternalError(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteInternalError(rr)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "expected status internal server error")
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), rr.Body.String(), "unexpected body content")
}

// TestWriteErrorResponse tests the WriteErrorResponse function.
func TestWriteErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	errors := []string{"error1", "error2"}

	WriteErrorResponse(rr, http.StatusBadRequest, errors)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "unexpected status code")

	var resp ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, errors, resp.Errors, "unexpected errors in response")
}

// TestWriteAPIResponse tests the WriteAPIResponse function.
func TestWriteAPIResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteAPIResponse(rr, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rr.Code, "unexpected status code")

	expectedBytes, err := json.MarshalIndent(Response{Data: data}, "", "  ")
	require.NoError(t, err)

	expected := string(expectedBytes)
	actual := rr.Body.String()

	assert.JSONEq(t, expected, actual, "unexpected data in response")
}

// TestHealthResponse tests the serialization of HealthResponse.
func TestHealthResponse(t *testing.T) {
	health := HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	bytes, err := json.Marshal(health)
	require.NoError(t, err)

	expected := `{"status":"ok","version":"1.0.0"}`
	assert.JSONEq(t, expected, string(bytes), "unexpected JSON serialization")
}

// TestGameResultResponse tests the serialization of GameResultResponse.
func TestGameResultResponse(t *testing.T) {
	id := uuid.New()
	gameResult := GameResultResponse{
		ID:                1,
		UserID:            id,
		GameStatus:        entity.GameStatusWin,
		TransactionSource: entity.TransactionSourceGame,
		TransactionID:     "txn123",
		Amount:            100.50,
		CreatedAt:         time.Now(),
	}

	bytes, err := json.Marshal(gameResult)
	require.NoError(t, err)

	var result GameResultResponse
	err = json.Unmarshal(bytes, &result)
	require.NoError(t, err)

	assert.Equal(t, gameResult.ID, result.ID, "unexpected ID")
}
