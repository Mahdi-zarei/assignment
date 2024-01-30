package test

import (
	db "assignment/db/utils"
	"assignment/pkg/common"
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func populateDBWithScenario(svc *ServiceTest) []types.GiftCardData {
	ctx := context.Background()
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()
	giftID1 := uuid.New()
	giftID2 := uuid.New()

	scenarios := []types.GiftCardData{
		{
			ID:           giftID1,
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           giftID2,
			GifterID:     userID3,
			GifteeID:     userID1,
			Status:       types.GiftCardStatusAccepted,
			IssueDate:    time.Now(),
			ResponseDate: time.Time{},
		},
	}

	for _, scenario := range scenarios {
		common.Must1(svc.giftCardRepo.InsertNew(ctx, &scenario))
	}

	return scenarios
}

func populateDBForFilter(svc *ServiceTest) []uuid.UUID {
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()
	IDs := []uuid.UUID{userID1, userID2, userID3}

	// userID1 got 2: 2 resp, 0 acc, 0 rej. sent 5: 1 resp, 3 acc, 1 rej.
	// userID2 got 4: 0 resp, 3 acc, 1 rej. sent 6: 5 resp, 0 acc, 1 rej.
	// userID3 got 5: 4 resp, 0 acc, 1 rej. sent 0: 0 resp, 0 acc, 0 rej.
	giftCards := []types.GiftCardData{
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusAccepted,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusAccepted,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusRejected,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusAccepted,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID1,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID1,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusRejected,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			GifterID:     userID2,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
	}

	for _, giftCard := range giftCards {
		common.Must1(svc.giftCardRepo.InsertNew(ctx, &giftCard))
	}

	return IDs
}

func TestInsertNew(t *testing.T) {
	svc := NewServiceTest()

	ctx := context.Background()
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()
	dupID := uuid.New()

	scenarios := []types.GiftCardData{
		{
			ID:           uuid.New(),
			GifterID:     userID1,
			GifteeID:     userID2,
			Status:       types.GiftCardStatusWaitingResponse,
			IssueDate:    time.Time{},
			ResponseDate: time.Time{},
		},
		{
			ID:           dupID,
			GifterID:     userID3,
			GifteeID:     userID1,
			Status:       types.GiftCardStatusAccepted,
			IssueDate:    time.Now(),
			ResponseDate: time.Time{},
		},
	}

	DupScenarios := []types.GiftCardData{
		{
			ID:           dupID,
			GifterID:     userID2,
			GifteeID:     userID3,
			Status:       types.GiftCardStatusRejected,
			IssueDate:    time.Now(),
			ResponseDate: time.Time{},
		},
	}

	for idx, scenario := range scenarios {
		t.Run("scenario "+strconv.Itoa(idx), func(t *testing.T) {
			scenarioCopy := scenario
			err := svc.giftCardRepo.InsertNew(ctx, &scenario)
			assert.Nil(t, err)
			assert.True(t, scenario.Equals(scenarioCopy))
		})
	}

	for idx, scenario := range DupScenarios {
		t.Run("dupScenario "+strconv.Itoa(idx), func(t *testing.T) {
			scenarioCopy := scenario
			err := svc.giftCardRepo.InsertNew(ctx, &scenario)
			assert.NotNil(t, err)
			assert.True(t, scenario.Equals(scenarioCopy))
		})
	}
}

func TestGetGiftData(t *testing.T) {
	svc := NewServiceTest()
	scenarios := populateDBWithScenario(svc)
	ctx := context.Background()

	for idx, scenario := range scenarios {
		t.Run("scenario "+strconv.Itoa(idx), func(t *testing.T) {
			res, err := svc.giftCardRepo.GetGiftData(ctx, scenario.ID)
			assert.Nil(t, err)
			if scenario.IssueDate.IsZero() { // handle the case where issue date is not provided
				assert.True(t, res.IssueDate.After(time.Now().Add(-1*time.Minute))) // make sure the issue date is in the near past
				scenario.IssueDate = res.IssueDate
			}
			assert.True(t, scenario.Equals(*res))
		})
	}
}

func TestUpdateGiftStatus(t *testing.T) {
	svc := NewServiceTest()
	scenarios := populateDBWithScenario(svc)
	ctx := context.Background()

	for idx, scenario := range scenarios {
		t.Run("scenario "+strconv.Itoa(idx), func(t *testing.T) {
			targetStatus := types.GiftCardStatusWaitingResponse
			err := svc.giftCardRepo.UpdateGiftStatus(ctx, scenario.ID, targetStatus)
			assert.Nil(t, err)

			res, err := svc.giftCardRepo.GetGiftData(ctx, scenario.ID)
			assert.Nil(t, err)
			assert.EqualValues(t, targetStatus, res.Status)
		})
	}
}

func TestGetGiftsByGifterID(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()
	fullPagination := db.PaginationData{
		PageNumber: 0,
		PageSize:   10000,
	}

	IDs := populateDBForFilter(svc)
	userID1 := IDs[0]
	userID2 := IDs[1]
	userID3 := IDs[2]

	res1, err := svc.giftCardRepo.GetGiftsByGifterID(ctx, userID1, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 5, len(res1))
	for _, res := range res1 {
		assert.EqualValues(t, userID1, res.GifterID)
	}

	res2, err := svc.giftCardRepo.GetGiftsByGifterID(ctx, userID2, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 6, len(res2))
	for _, res := range res2 {
		assert.EqualValues(t, userID2, res.GifterID)
	}

	res3, err := svc.giftCardRepo.GetGiftsByGifterID(ctx, userID3, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(res3))

	// test correct pagination for this function
	pageRes, err := svc.giftCardRepo.GetGiftsByGifterID(ctx, userID2, db.PaginationData{
		PageNumber: 0,
		PageSize:   2,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(pageRes))

	pageRes, err = svc.giftCardRepo.GetGiftsByGifterID(ctx, userID2, db.PaginationData{
		PageNumber: 3,
		PageSize:   2,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(pageRes))
}

func TestGetGiftsByGifteeID(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()
	fullPagination := db.PaginationData{
		PageNumber: 0,
		PageSize:   10000,
	}

	IDs := populateDBForFilter(svc)
	userID1 := IDs[0]
	userID2 := IDs[1]
	userID3 := IDs[2]

	res1, err := svc.giftCardRepo.GetGiftsByGifteeID(ctx, userID1, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(res1))
	for _, res := range res1 {
		assert.EqualValues(t, userID1, res.GifteeID)
	}

	res2, err := svc.giftCardRepo.GetGiftsByGifteeID(ctx, userID2, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 4, len(res2))
	for _, res := range res2 {
		assert.EqualValues(t, userID2, res.GifteeID)
	}

	res3, err := svc.giftCardRepo.GetGiftsByGifteeID(ctx, userID3, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 5, len(res3))
	for _, res := range res3 {
		assert.EqualValues(t, userID3, res.GifteeID)
	}

	// test pagination
	pageRes, err := svc.giftCardRepo.GetGiftsByGifteeID(ctx, userID1, db.PaginationData{
		PageNumber: 1,
		PageSize:   1,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(pageRes))

	pageRes, err = svc.giftCardRepo.GetGiftsByGifteeID(ctx, userID1, db.PaginationData{
		PageNumber: 10,
		PageSize:   1,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(pageRes))
}

func TestGetGiftsByGifterIDAndStatus(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()
	fullPagination := db.PaginationData{
		PageNumber: 0,
		PageSize:   10000,
	}

	IDs := populateDBForFilter(svc)
	userID1 := IDs[0]
	userID2 := IDs[1]
	userID3 := IDs[2]

	res1Waiting, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID1, types.GiftCardStatusWaitingResponse, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(res1Waiting))
	for _, res := range res1Waiting {
		assert.EqualValues(t, userID1, res.GifterID)
		assert.EqualValues(t, types.GiftCardStatusWaitingResponse, res.Status)
	}

	res1Accepted, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID1, types.GiftCardStatusAccepted, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(res1Accepted))
	for _, res := range res1Accepted {
		assert.EqualValues(t, userID1, res.GifterID)
		assert.EqualValues(t, types.GiftCardStatusAccepted, res.Status)
	}

	res2Accepted, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID2, types.GiftCardStatusAccepted, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(res2Accepted))
	for _, res := range res2Accepted {
		assert.EqualValues(t, userID2, res.GifterID)
		assert.EqualValues(t, types.GiftCardStatusAccepted, res.Status)
	}

	res2Rejected, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID2, types.GiftCardStatusRejected, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(res2Rejected))
	for _, res := range res2Rejected {
		assert.EqualValues(t, userID2, res.GifterID)
		assert.EqualValues(t, types.GiftCardStatusRejected, res.Status)
	}

	res3Waiting, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID3, types.GiftCardStatusWaitingResponse, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(res3Waiting))
	for _, res := range res3Waiting {
		assert.EqualValues(t, userID3, res.GifterID)
		assert.EqualValues(t, types.GiftCardStatusWaitingResponse, res.Status)
	}

	// test pagination
	pageRes2Waiting, err := svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID2, types.GiftCardStatusWaitingResponse, db.PaginationData{
		PageNumber: 0,
		PageSize:   3,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(pageRes2Waiting))

	pageRes2Waiting, err = svc.giftCardRepo.GetGiftsByGifterIDAndStatus(ctx, userID2, types.GiftCardStatusWaitingResponse, db.PaginationData{
		PageNumber: 1,
		PageSize:   3,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(pageRes2Waiting))
}

func TestGetGiftsByGifteeIDAndStatus(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()
	fullPagination := db.PaginationData{
		PageNumber: 0,
		PageSize:   10000,
	}

	IDs := populateDBForFilter(svc)
	//userID1 := IDs[0]
	userID2 := IDs[1]
	userID3 := IDs[2]

	res2Waiting, err := svc.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, userID2, types.GiftCardStatusWaitingResponse, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(res2Waiting))
	for _, res := range res2Waiting {
		assert.EqualValues(t, userID2, res.GifteeID)
		assert.EqualValues(t, types.GiftCardStatusWaitingResponse, res.Status)
	}

	res2Accepted, err := svc.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, userID2, types.GiftCardStatusAccepted, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(res2Accepted))
	for _, res := range res2Accepted {
		assert.EqualValues(t, userID2, res.GifteeID)
		assert.EqualValues(t, types.GiftCardStatusAccepted, res.Status)
	}

	res3Rejected, err := svc.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, userID3, types.GiftCardStatusRejected, fullPagination)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(res3Rejected))
	for _, res := range res3Rejected {
		assert.EqualValues(t, userID3, res.GifteeID)
		assert.EqualValues(t, types.GiftCardStatusRejected, res.Status)
	}

	// test pagination
	pageRes3Waiting, err := svc.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, userID3, types.GiftCardStatusWaitingResponse, db.PaginationData{
		PageNumber: 1,
		PageSize:   2,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(pageRes3Waiting))

	pageRes3Waiting, err = svc.giftCardRepo.GetGiftsByGifteeIDAndStatus(ctx, userID3, types.GiftCardStatusWaitingResponse, db.PaginationData{
		PageNumber: 2,
		PageSize:   2,
	})
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(pageRes3Waiting))
}
