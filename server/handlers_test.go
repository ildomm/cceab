package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ildomm/cceab/entity"
	"github.com/ildomm/cceab/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHealthHandlerSuccess tests the Health function for a successful response.
func TestHealthHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server := Server{} // Assuming Server struct exists
		server.HealthHandler(w, r)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := Response{}
	expected.Data = HealthResponse{Status: "pass", Version: "v1"}
	body, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	var actual Response
	actual.Data = HealthResponse{}

	err = json.Unmarshal(body, &actual)
	require.NoError(t, err)
}

// TestGameResultFuncSuccess tests the CreateGameResultFunc for a successful response using a real server.
func TestGameResultFuncSuccess(t *testing.T) {
	mockDAO := test_helpers.NewMockGameResultDAO()

	// Set up mock expectations
	testGameResult := &entity.GameResult{
		ID:            1,
		UserID:        uuid.New(),
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
		CreatedAt:     time.Now(),
	}
	mockDAO.On("CreateGameResult",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(testGameResult, nil)

	// Create the server and set the mock manager
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithGameResultManager(mockDAO)
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, testGameResult.UserID.String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "CreateDeviceFunc returned wrong status code")

	// Decode the response
	var respBody Response
	respBody.Data = GameResultResponse{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
}

// TestGameResultFuncInvalidRequest tests the CreateGameResultFunc with an invalid request body.
func TestGameResultFuncInvalidRequest(t *testing.T) {
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request with invalid body
	body := []byte(`not a json`)
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, uuid.New().String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "CreateDeviceFunc returned wrong status code for invalid request body")
}

// TestGameResultFuncInvalidUserID tests the CreateGameResultFunc with an invalid user ID.
func TestGameResultFuncInvalidUserID(t *testing.T) {
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request with invalid user ID
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/invalid-user-id/game_results", port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "CreateDeviceFunc returned wrong status code for invalid user ID")
}

// TestGameResultFuncUserNotFound tests the CreateGameResultFunc when the user is not found.
func TestGameResultFuncUserNotFound(t *testing.T) {
	mockDAO := test_helpers.NewMockGameResultDAO()

	// Set up mock expectations
	mockDAO.On("CreateGameResult",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, entity.ErrUserNotFound)

	// Create the server and set the mock manager
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithGameResultManager(mockDAO)
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, uuid.New().String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "CreateDeviceFunc returned wrong status code for user not found")
}

// TestGameResultFuncTransactionIDExists tests the CreateGameResultFunc when the transaction ID already exists.
func TestGameResultFuncTransactionIDExists(t *testing.T) {
	mockDAO := test_helpers.NewMockGameResultDAO()

	// Set up mock expectations
	mockDAO.On("CreateGameResult",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, entity.ErrTransactionIdExists)

	// Create the server and set the mock manager
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithGameResultManager(mockDAO)
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, uuid.New().String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode, "CreateDeviceFunc returned wrong status code for transaction ID exists")
}

// TestGameResultFuncUserNegativeBalance tests the CreateGameResultFunc when the user has a negative balance.
func TestGameResultFuncUserNegativeBalance(t *testing.T) {
	mockDAO := test_helpers.NewMockGameResultDAO()

	// Set up mock expectations
	mockDAO.On("CreateGameResult",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, entity.ErrUserNegativeBalance)

	// Create the server and set the mock manager
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithGameResultManager(mockDAO)
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "win",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, uuid.New().String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode, "CreateDeviceFunc returned wrong status code for user with negative balance")
}

// TestGameResultFuncInvalidGameStatus tests the CreateGameResultFunc with an invalid game status.
func TestGameResultFuncInvalidGameStatus(t *testing.T) {
	mockDAO := test_helpers.NewMockGameResultDAO()

	// Set up mock expectations
	mockDAO.On("CreateGameResult",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, entity.ErrInvalidGameStatus)

	// Create the server and set the mock manager
	server := NewServer()
	port := rand.Intn(1000) + 8000
	server.WithGameResultManager(mockDAO)
	server.WithListenAddress(port)

	go func() {
		err := server.Run()
		assert.NoError(t, err, "server failed to run")
	}()

	// Create the request body
	reqBody := CreateGameResultRequest{
		GameStatus:    "invalid-status",
		Amount:        100,
		TransactionID: "123",
	}
	body, _ := json.Marshal(reqBody)

	// Create the request
	url := fmt.Sprintf("http://localhost:%d/api/v1/users/%s/game_results", port, uuid.New().String())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("source-type", string(entity.TransactionSourceGame))

	// Use httptest to create a server
	testServer := httptest.NewServer(server.router())
	defer testServer.Close()

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request to server failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "CreateDeviceFunc returned wrong status code for invalid game status")
}
