package giftcard

import (
	db "assignment/db/utils"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
)

type GiftCardModule interface {
	IssueNewGiftCard(ctx context.Context, gifterID uuid.UUID, gifteeID uuid.UUID) (*types.GiftCardData, error)

	RespondToGift(ctx context.Context, giftID uuid.UUID, targetStatus types.GiftCardStatus) error

	InquireGiftCard(ctx context.Context, giftID uuid.UUID) (*types.GiftCardData, error)
	GetListOfSentGiftCards(ctx context.Context, gifterID uuid.UUID, wantedStatus *types.GiftCardStatus, pagination db.PaginationData) ([]*types.GiftCardData, error)
	GetListOfReceivedGiftCards(ctx context.Context, gifteeID uuid.UUID, wantedStatus *types.GiftCardStatus, pagination db.PaginationData) ([]*types.GiftCardData, error)
}
