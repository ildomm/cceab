package entity

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserNegativeBalance = errors.New("negative balance not allowed")
var ErrGameResultNotFound = errors.New("game result not found")
var ErrTransactionIdExists = errors.New("transaction id already exists")
var ErrInvalidGameStatus = errors.New("invalid game status")
var ErrInvalidTransactionSource = errors.New("invalid transaction source")
var ErrInvalidTransactionStatus = errors.New("invalid transaction status")
