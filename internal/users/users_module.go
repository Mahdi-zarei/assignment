package users

import (
	"assignment/db/users_repo"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type UsersModuleImpl struct {
	usersRepo users_repo.UsersRepo
	logger    *logrus.Logger
}

func NewUsersModule(repo users_repo.UsersRepo, logger *logrus.Logger) UsersModule {
	return &UsersModuleImpl{
		usersRepo: repo,
		logger:    logger,
	}
}

func (u *UsersModuleImpl) RegisterNewUser(ctx context.Context, name string) (uuid.UUID, error) {
	const spot = "RegisterNewUser"

	userID := uuid.New()
	userData := types.UserData{
		ID:           userID,
		Name:         name,
		RegisterDate: time.Now(),
	}

	err := u.usersRepo.InsertNew(ctx, &userData)
	if err != nil {
		u.logger.Errorf("[%s] Failed to insert new user: %s", spot, err)
		return uuid.Nil, err
	}

	return userID, nil
}

func (u *UsersModuleImpl) GetUserData(ctx context.Context, userID uuid.UUID) (*types.UserData, error) {
	const spot = "GetUserData"

	res, err := u.usersRepo.GetUserData(ctx, userID)
	if err != nil {
		u.logger.Errorf("[%s] Failed to get user data: %s", spot, err)
		return nil, err
	}

	return res, nil
}
