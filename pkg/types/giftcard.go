package types

import (
	"github.com/google/uuid"
	"time"
)

type GiftCardStatus int32

const (
	GiftCardStatusUnknown GiftCardStatus = iota
	GiftCardStatusWaitingResponse
	GiftCardStatusAccepted
	GiftCardStatusRejected
)

type GiftCardData struct {
	ID           uuid.UUID      `json:"id"`
	GifterID     uuid.UUID      `json:"gifter_id"`
	GifteeID     uuid.UUID      `json:"giftee_id"`
	Status       GiftCardStatus `json:"status"`
	IssueDate    time.Time      `json:"issue_date"`
	ResponseDate time.Time      `json:"response_date"`
}

func (g *GiftCardData) Equals(target GiftCardData) bool {
	return g.ID == target.ID &&
		g.GifterID == target.GifterID &&
		g.GifteeID == target.GifteeID &&
		g.Status == target.Status &&
		g.IssueDate.UnixMilli() == target.IssueDate.UnixMilli() &&
		g.ResponseDate.UnixMilli() == target.ResponseDate.UnixMilli()
}
