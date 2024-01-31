package users

import (
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
)

type UsersModule interface {
	RegisterNewUser(ctx context.Context, name string) (*types.UserData, error)
	GetUserData(ctx context.Context, userID uuid.UUID) (*types.UserData, error)
}
