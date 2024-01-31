package api

import "time"

const (
	UserID         = "user-id"
	UserName       = "username"
	GiftID         = "gift-id"
	GifterID       = "gifter-id"
	GifteeID       = "giftee-id"
	ResponseToGift = "response-to-gift"
	WantedStatus   = "wanted-status"
	PageNumber     = "page-number"
	PageSize       = "page-size"

	UnknownStatus  = "unknown"
	PendingStatus  = "pending"
	AcceptedStatus = "accepted"
	RejectedStatus = "rejected"

	RequestTimeout = 10 * time.Second
)
