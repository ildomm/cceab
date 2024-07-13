package dao

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserNegativeBalance = errors.New("negative balance not allowed")

type UserDAO interface {
	UpdateUserBalance(ctx context.Context, userId uuid.UUID, amount float64, validationStatus bool) error
}
