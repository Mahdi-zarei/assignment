package users_repo

import (
	"assignment/pkg/common"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoImpl struct {
	db *pgxpool.Pool
}

func NewUsersRepo(ctx context.Context, db *pgxpool.Pool) UsersRepo {
	common.Must2(db.Exec(ctx, `CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    name TEXT,
    register_date TIMESTAMPTZ DEFAULT NOW()
)`))

	return &UserRepoImpl{
		db: db,
	}
}

func (u *UserRepoImpl) InsertNew(ctx context.Context, userData *types.UserData) error {
	_, err := u.db.Exec(ctx, `INSERT INTO users(id, name, register_date) VALUES ($1,$2,$3)`,
		userData.ID,
		userData.Name,
		common.NowIfZero(userData.RegisterDate))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoImpl) GetUserData(ctx context.Context, userID uuid.UUID) (*types.UserData, error) {
	var res types.UserData

	err := u.db.QueryRow(ctx, `SELECT id, name, register_date FROM users WHERE id=$1`, userID).Scan(
		&res.ID,
		&res.Name,
		&res.RegisterDate)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
