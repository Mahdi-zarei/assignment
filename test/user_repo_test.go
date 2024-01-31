package test

import (
	"assignment/pkg/types"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestInsertUser(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()

	dupID := uuid.New()

	scenarios := []types.UserData{
		{
			ID:           uuid.New(),
			Name:         "Mahdi",
			RegisterDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			Name:         "Ali",
			RegisterDate: time.Now(),
		},
		{
			ID:           dupID,
			Name:         "Dup",
			RegisterDate: time.Now(),
		},
	}

	failScenarios := []types.UserData{
		{
			ID:           dupID,
			Name:         "Asdf",
			RegisterDate: time.Now(),
		},
	}

	for idx, scenario := range scenarios {
		t.Run("scenario "+strconv.Itoa(idx), func(t *testing.T) {
			err := svc.usersRepo.InsertNew(ctx, &scenario)
			assert.Nil(t, err)
		})
	}

	for idx, scenario := range failScenarios {
		t.Run("failScenario "+strconv.Itoa(idx), func(t *testing.T) {
			err := svc.usersRepo.InsertNew(ctx, &scenario)
			assert.NotNil(t, err)
		})
	}
}

func TestGetUserData(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()

	users := []types.UserData{
		{
			ID:           uuid.New(),
			Name:         "Mahdi",
			RegisterDate: time.Time{},
		},
		{
			ID:           uuid.New(),
			Name:         "Ali",
			RegisterDate: time.Now(),
		},
	}

	for idx, user := range users {
		t.Run("scenario "+strconv.Itoa(idx), func(t *testing.T) {
			err := svc.usersRepo.InsertNew(ctx, &user)
			assert.Nil(t, err)

			res, err := svc.usersRepo.GetUserData(ctx, user.ID)
			assert.Nil(t, err)
			if user.RegisterDate.IsZero() {
				assert.True(t, res.RegisterDate.After(time.Now().Add(-1*time.Minute))) // assert register date has been properly overriden
				user.RegisterDate = res.RegisterDate
			}
			assert.True(t, user.Equals(*res))
		})
	}
}
