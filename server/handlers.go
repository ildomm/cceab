package server

import (
	"github.com/ildomm/cceab/dao"
	"net/http"
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
	userDAO       dao.UserDAO
	gameResultDAO dao.GameResultDAO
}

func NewGameResultHandler(userDAO dao.UserDAO, gameResultDAO dao.GameResultDAO) *gameResultHandler {
	return &gameResultHandler{
		userDAO:       userDAO,
		gameResultDAO: gameResultDAO,
	}
}

// CreateGameResultFunc handles the request to create a new game result.
func (h *gameResultHandler) CreateGameResultFunc(w http.ResponseWriter, r *http.Request) {
	// TODO
	WriteAPIResponse(w, http.StatusCreated, "")
}
