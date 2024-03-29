package types

import (
	"github.com/google/uuid"
	"time"
)

// since handling user related aspects is not a concern, we will define a user with the minimal data needed

type UserData struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RegisterDate time.Time `json:"register_date"`
}

func (u *UserData) Equals(target UserData) bool {
	return u.ID == target.ID &&
		u.Name == target.Name &&
		u.RegisterDate.UnixMilli() == target.RegisterDate.UnixMilli()
}
