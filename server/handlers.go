package server

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ildomm/cceab/dao"
	"github.com/ildomm/cceab/entity"
	"net/http"
	"strconv"
	"strings"
)

// HealthHandler evaluates the health of the service and writes a standardized response.
func (s *Server) HealthHandler(response http.ResponseWriter, request *http.Request) {
	health := HealthResponse{
		Status:  "pass",
		Version: "v1",
	}

	WriteAPIResponse(response, http.StatusOK, health)
}

// gameResultHandler handles all requests related to game results.
type gameResultHandler struct {
	gameResultDAO dao.GameResultDAO
}

func NewGameResultHandler(gameResultDAO dao.GameResultDAO) *gameResultHandler {
	return &gameResultHandler{
		gameResultDAO: gameResultDAO,
	}
}

// CreateGameResultFunc handles the request to create a new game result.
func (h *gameResultHandler) CreateGameResultFunc(w http.ResponseWriter, r *http.Request) {
	// Validate the headers.
	transactionSource := entity.ParseTransactionSource(strings.ToLower(r.Header.Get("Source-Type")))
	if transactionSource == nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{entity.ErrInvalidTransactionSource.Error()})
		return

	}

	// Validate the request body.
	var req CreateGameResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{entity.ErrRequestPayload.Error()})
		return
	}

	// Validate amount type cast
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{entity.ErrInvalidAmount.Error()})
		return
	}

	// Validate the game status.
	if req.GameStatus != entity.GameStatusWin && req.GameStatus != entity.GameStatusLost {
		WriteErrorResponse(w, http.StatusBadRequest, []string{entity.ErrInvalidGameStatus.Error()})
	}

	// Extract the user ID from the request path.
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{entity.ErrInvalidUser.Error()})
		return
	}

	// Perform the business logic.
	gameResult, err := h.gameResultDAO.CreateGameResult(r.Context(), userId, req.GameStatus, amount, *transactionSource, req.TransactionID)
	if err != nil {

		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})

		case errors.Is(err, entity.ErrTransactionIdExists) || errors.Is(err, entity.ErrUserNegativeBalance):
			WriteErrorResponse(w, http.StatusNotAcceptable, []string{err.Error()})

		default:
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}

		return
	}

	gameResultResponse := transformGameResultResponse(*gameResult)
	WriteAPIResponse(w, http.StatusCreated, gameResultResponse)
}

// Transform entity.GameResult to server.GameResultResponse
func transformGameResultResponse(gameResult entity.GameResult) GameResultResponse {
	return GameResultResponse{
		ID:                gameResult.ID,
		UserID:            gameResult.UserID,
		GameStatus:        gameResult.GameStatus,
		Amount:            gameResult.Amount,
		TransactionSource: gameResult.TransactionSource,
		TransactionID:     gameResult.TransactionID,
		CreatedAt:         gameResult.CreatedAt,
	}
}
