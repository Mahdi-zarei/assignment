package users_repo

import (
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
)

// UsersRepo defines the function scheme for users repository
// since users and their management are not of concern here, we define a minimalistic db
type UsersRepo interface {
	InsertNew(ctx context.Context, userData *types.UserData) error

	GetUserData(ctx context.Context, userID uuid.UUID) (*types.UserData, error)
}
