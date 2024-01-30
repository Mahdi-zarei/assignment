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
