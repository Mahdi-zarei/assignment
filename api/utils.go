package api

import (
	db "assignment/db/utils"
	"assignment/pkg/common"
	"assignment/pkg/types"
	"net/http"
	"strconv"
)

func writeError(response http.ResponseWriter, err error, status int) {
	response.WriteHeader(status)
	response.Write([]byte(err.Error()))
}

func parseStatus(status string) (*types.GiftCardStatus, error) {
	switch status {
	case "":
		return nil, nil
	case UnknownStatus:
		return common.PtrTo(types.GiftCardStatusUnknown), nil
	case PendingStatus:
		return common.PtrTo(types.GiftCardStatusWaitingResponse), nil
	case AcceptedStatus:
		return common.PtrTo(types.GiftCardStatusAccepted), nil
	case RejectedStatus:
		return common.PtrTo(types.GiftCardStatusRejected), nil
	default:
		return nil, types.ErrInvalidStatus
	}
}

func parseResponse(response string) (types.GiftCardStatus, error) {
	switch response {
	case AcceptedStatus:
		return types.GiftCardStatusAccepted, nil
	case RejectedStatus:
		return types.GiftCardStatusRejected, nil
	default:
		return types.GiftCardStatusUnknown, types.ErrInvalidResponse
	}
}

func parsePagination(pageSize string, pageNumber string) (*db.PaginationData, error) {
	if pageSize == "" || pageNumber == "" {
		return &db.DefaultPagination, nil
	}

	paginationData := db.PaginationData{
		PageNumber: int32(common.GetVal(strconv.Atoi(pageNumber))),
		PageSize:   int32(common.GetVal(strconv.Atoi(pageSize))),
	}
	if paginationData.IsValid() {
		return &paginationData, nil
	}

	return nil, types.ErrInvalidPagination
}
