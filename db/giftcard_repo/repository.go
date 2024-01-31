package giftcard_repo

import (
	db "assignment/db/utils"
	"assignment/pkg/common"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GiftCardRepoImpl struct {
	db *pgxpool.Pool
}

func NewGiftCardRepo(ctx context.Context, db *pgxpool.Pool) GiftCardRepo {
	common.Must2(db.Exec(ctx, `CREATE TABLE IF NOT EXISTS gift_card(
    id UUID PRIMARY KEY,
    gifter_id UUID NOT NULL,
    giftee_id UUID NOT NULL,
    status INTEGER,
    issue_date TIMESTAMPTZ DEFAULT NOW(),
    response_date TIMESTAMPTZ
)`))

	common.Must2(db.Exec(ctx, `CREATE INDEX IF NOT EXISTS status_idx ON gift_card(status)`))
	common.Must2(db.Exec(ctx, `CREATE INDEX IF NOT EXISTS gifter_idx ON gift_card(gifter_id)`))
	common.Must2(db.Exec(ctx, `CREATE INDEX IF NOT EXISTS giftee_idx ON gift_card(giftee_id)`))

	return &GiftCardRepoImpl{
		db: db,
	}
}

func (g *GiftCardRepoImpl) InsertNew(ctx context.Context, giftData *types.GiftCardData) error {
	_, err := g.db.Exec(ctx, "INSERT INTO gift_card(id, gifter_id, giftee_id, status, issue_date, response_date) VALUES($1,$2,$3,$4,$5,$6)",
		giftData.ID,
		giftData.GifterID,
		giftData.GifteeID,
		giftData.Status,
		common.NowIfZero(giftData.IssueDate),
		giftData.ResponseDate)
	if err != nil {
		return err
	}

	return nil
}

func (g *GiftCardRepoImpl) UpdateGiftStatus(ctx context.Context, giftID uuid.UUID, targetStatus types.GiftCardStatus, setResponseDate bool) error {
	query := ""
	if setResponseDate {
		query = "UPDATE gift_card SET status = $1, response_date=NOW() WHERE id = $2"
	} else {
		query = "UPDATE gift_card SET status = $1 WHERE id = $2"
	}

	res, err := g.db.Exec(ctx, query, targetStatus, giftID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return types.ErrNotFound
	}

	return nil
}

func (g *GiftCardRepoImpl) GetGiftData(ctx context.Context, id uuid.UUID) (*types.GiftCardData, error) {
	var res types.GiftCardData

	err := g.db.QueryRow(ctx, "SELECT gifter_id, giftee_id, status, issue_date, response_date FROM gift_card WHERE id = $1", id).Scan(
		&res.GifterID,
		&res.GifteeID,
		&res.Status,
		&res.IssueDate,
		&res.ResponseDate)
	if err != nil {
		return nil, err
	}

	res.ID = id
	return &res, nil
}

func (g *GiftCardRepoImpl) GetGiftsByGifterID(ctx context.Context, gifterID uuid.UUID, paginationData db.PaginationData) ([]*types.GiftCardData, error) {
	var res []*types.GiftCardData
	limit, offset := paginationData.GetLimitAndOffset()

	rows, err := g.db.Query(ctx, "SELECT id, gifter_id, giftee_id, status, issue_date, response_date FROM gift_card "+
		"WHERE gifter_id = $1 LIMIT $2 OFFSET $3", gifterID, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rec types.GiftCardData
		err = rows.Scan(
			&rec.ID,
			&rec.GifterID,
			&rec.GifteeID,
			&rec.Status,
			&rec.IssueDate,
			&rec.ResponseDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &rec)
	}

	return res, nil
}

func (g *GiftCardRepoImpl) GetGiftsByGifteeID(ctx context.Context, gifteeID uuid.UUID, paginationData db.PaginationData) ([]*types.GiftCardData, error) {
	var res []*types.GiftCardData
	limit, offset := paginationData.GetLimitAndOffset()

	rows, err := g.db.Query(ctx, "SELECT id, gifter_id, giftee_id, status, issue_date, response_date FROM gift_card "+
		"WHERE giftee_id = $1 LIMIT $2 OFFSET $3", gifteeID, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rec types.GiftCardData
		err = rows.Scan(
			&rec.ID,
			&rec.GifterID,
			&rec.GifteeID,
			&rec.Status,
			&rec.IssueDate,
			&rec.ResponseDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &rec)
	}

	return res, nil
}

func (g *GiftCardRepoImpl) GetGiftsByGifterIDAndStatus(ctx context.Context, gifterID uuid.UUID, wantedStatus types.GiftCardStatus, paginationData db.PaginationData) ([]*types.GiftCardData, error) {
	var res []*types.GiftCardData
	limit, offset := paginationData.GetLimitAndOffset()

	rows, err := g.db.Query(ctx, "SELECT id, gifter_id, giftee_id, status, issue_date, response_date FROM gift_card "+
		"WHERE gifter_id = $1 AND status = $2 LIMIT $3 OFFSET $4", gifterID, wantedStatus, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rec types.GiftCardData
		err = rows.Scan(
			&rec.ID,
			&rec.GifterID,
			&rec.GifteeID,
			&rec.Status,
			&rec.IssueDate,
			&rec.ResponseDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &rec)
	}

	return res, nil
}

func (g *GiftCardRepoImpl) GetGiftsByGifteeIDAndStatus(ctx context.Context, gifteeID uuid.UUID, wantedStatus types.GiftCardStatus, paginationData db.PaginationData) ([]*types.GiftCardData, error) {
	var res []*types.GiftCardData
	limit, offset := paginationData.GetLimitAndOffset()

	rows, err := g.db.Query(ctx, "SELECT id, gifter_id, giftee_id, status, issue_date, response_date FROM gift_card "+
		"WHERE giftee_id = $1 AND status = $2 LIMIT $3 OFFSET $4", gifteeID, wantedStatus, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rec types.GiftCardData
		err = rows.Scan(
			&rec.ID,
			&rec.GifterID,
			&rec.GifteeID,
			&rec.Status,
			&rec.IssueDate,
			&rec.ResponseDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &rec)
	}

	return res, nil
}
