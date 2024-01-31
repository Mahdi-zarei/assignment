package giftcard

import (
	"assignment/db/giftcard_repo"
	"assignment/db/users_repo"
	db "assignment/db/utils"
	"assignment/pkg/types"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type GiftCardModuleImpl struct {
	giftCardRepo giftcard_repo.GiftCardRepo
	usersRepo    users_repo.UsersRepo
	logger       *logrus.Logger
}

func NewGiftCardModule(repo giftcard_repo.GiftCardRepo, usersRepo users_repo.UsersRepo, logger *logrus.Logger) GiftCardModule {
	return &GiftCardModuleImpl{
		giftCardRepo: repo,
		usersRepo:    usersRepo,
		logger:       logger,
	}
}

func (g *GiftCardModuleImpl) IssueNewGiftCard(ctx context.Context, gifterID uuid.UUID, gifteeID uuid.UUID) (*types.GiftCardData, error) {
	const spot = "IssueNewGiftCard"

	_, err := g.usersRepo.GetUserData(ctx, gifteeID)
	if err != nil {
		g.logger.Errorf("[%s] Failed to get giftee data: %s", spot, err)
		if errors.Is(err, pgx.ErrNoRows) {
			err = types.ErrGifteeInvalid
		}
		return nil, err
	}

	giftID := uuid.New()
	giftData := types.GiftCardData{
		ID:        giftID,
		GifterID:  gifterID,
		GifteeID:  gifteeID,
		Status:    types.GiftCardStatusWaitingResponse,
		IssueDate: time.Now(),
	}

	err = g.giftCardRepo.InsertNew(ctx, &giftData)
	if err != nil {
		g.logger.Errorf("[%s] Failed to insert new giftcard: %s", spot, err)
		return nil, err
	}

	return &giftData, nil
}

func (g *GiftCardModuleImpl) RespondToGift(ctx context.Context, giftID uuid.UUID, targetStatus types.GiftCardStatus) error {
	const spot = "RespondToGift"
	err := g.giftCardRepo.UpdateGiftStatus(ctx, giftID, targetStatus, true)
	if err != nil {
		g.logger.Errorf("[%s] Failed to update gift status: %s", spot, err)
		return err
	}

	return nil
}

func (g *GiftCardModuleImpl) InquireGiftCard(ctx context.Context, giftID uuid.UUID) (*types.GiftCardData, error) {
	const spot = "InquireGiftCard"
	res, err := g.giftCardRepo.GetGiftData(ctx, giftID)
	if err != nil {
		g.logger.Errorf("[%s] Failed to get gift data: %s", spot, err)
		return nil, err
	}

	return res, nil
}

func (g *GiftCardModuleImpl) GetListOfSentGiftCards(ctx context.Context, gifterID uuid.UUID, wantedStatus *types.GiftCardStatus, pagination db.PaginationData) ([]*types.GiftCardData, error) {
	const spot = "GetListOfSentGiftCards"
	var res []*types.GiftCardData
	var err error

	// practically we should check if the user exists here, but for the sake of simplicity we just hit the DB
	if wantedStatus == nil {
		res, err = g.giftCardRepo.GetGiftsByGifterID(ctx, gifterID, pagination)
	} else {
		res, err = g.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, gifterID, *wantedStatus, pagination)
	}
	if err != nil {
		g.logger.Errorf("[%s] Failed to get sent gift cards: %s", spot, err)
		return nil, err
	}

	return res, nil
}

func (g *GiftCardModuleImpl) GetListOfReceivedGiftCards(ctx context.Context, gifteeID uuid.UUID, wantedStatus *types.GiftCardStatus, pagination db.PaginationData) ([]*types.GiftCardData, error) {
	const spot = "GetListOfReceivedGiftCards"
	var res []*types.GiftCardData
	var err error

	// practically we should check if the user exists here, but for the sake of simplicity we just hit the DB
	if wantedStatus == nil {
		res, err = g.giftCardRepo.GetGiftsByGifteeID(ctx, gifteeID, pagination)
	} else {
		res, err = g.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, gifteeID, *wantedStatus, pagination)
	}
	if err != nil {
		g.logger.Errorf("[%s] Failed to get received gift cards: %s", spot, err)
		return nil, err
	}

	return res, nil
}
