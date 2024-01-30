package giftcard_repo

import (
	"assignment/db/utils"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
)

// GiftCardRepo Defines the function schemes for the GiftCard Repository
// since we want to implement pagination it's better to apply filters at the DB level, otherwise given that each user would've most likely
// had few gifts related to him, it could've been viable to apply the status filtering at the application level
type GiftCardRepo interface {
	InsertNew(ctx context.Context, giftData *types.GiftCardData) error

	UpdateGiftStatus(ctx context.Context, giftID uuid.UUID, targetStatus types.GiftCardStatus) error

	GetGiftData(ctx context.Context, id uuid.UUID) (*types.GiftCardData, error)
	GetGiftsByGifterID(ctx context.Context, gifterID uuid.UUID, paginationData db.PaginationData) ([]*types.GiftCardData, error)
	GetGiftsByGifteeID(ctx context.Context, gifteeID uuid.UUID, paginationData db.PaginationData) ([]*types.GiftCardData, error)
	GetGiftsByGifterIDAndStatus(ctx context.Context, gifterID uuid.UUID, wantedStatus types.GiftCardStatus, paginationData db.PaginationData) ([]*types.GiftCardData, error)
	GetGiftsByGifteeIDAndStatus(ctx context.Context, gifteeID uuid.UUID, wantedStatus types.GiftCardStatus, paginationData db.PaginationData) ([]*types.GiftCardData, error)
}
