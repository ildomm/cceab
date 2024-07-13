package server

import "github.com/ildomm/cceab/entity"

type CreateGameResultRequest struct {
	GameStatus    entity.GameStatus `json:"state"`
	Amount        float64           `json:"amount"`
	TransactionID string            `json:"transactionId"`
}
