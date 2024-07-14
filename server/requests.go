package server

import "github.com/ildomm/cceab/entity"

type CreateGameResultRequest struct {
	GameStatus    entity.GameStatus `json:"state"`
	Amount        string            `json:"amount"`
	TransactionID string            `json:"transactionId"`
}
