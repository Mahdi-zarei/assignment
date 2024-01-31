package api

import (
	"assignment/config"
	"assignment/internal/giftcard"
	"assignment/internal/users"
	"assignment/pkg/common"
	"assignment/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ApiHandler struct {
	giftCardModule giftcard.GiftCardModule
	usersModule    users.UsersModule
	logger         *logrus.Logger
}

func NewApiHandler(giftCardModule giftcard.GiftCardModule, usersModule users.UsersModule, logger *logrus.Logger, listenConf config.ListenConf) *ApiHandler {
	handler := &ApiHandler{
		giftCardModule: giftCardModule,
		usersModule:    usersModule,
		logger:         logger,
	}

	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/register-user", handler.RegisterUser)
	httpHandler.HandleFunc("/send-new-gift", handler.SendNewGift)
	httpHandler.HandleFunc("/respond-to-gift", handler.RespondToGift)
	httpHandler.HandleFunc("/inquire-sent-gifts", handler.InquireSentGifts)
	httpHandler.HandleFunc("/inquire-received-gifts", handler.InquireReceivedGifts)
	go func() {
		addr := fmt.Sprintf("%s:%d", listenConf.Host, listenConf.Port)
		logger.Fatal(http.ListenAndServe(addr, httpHandler))
	}()

	return handler
}

// RegisterUser registers a new user, it doesn't check anything and returns the entire user data if succeeds
func (h *ApiHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
	const spot = "api/RegisterUser"
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	username := request.URL.Query().Get(UserName)
	res, err := h.usersModule.RegisterNewUser(ctx, username)
	if err != nil {
		writeError(response, err, 500)
		return
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		h.logger.Errorf("[%s] Failed to marshal json: %s", spot, err)
		writeError(response, err, 500)
		return
	}

	response.Write(jsonResp)
}

// SendNewGift reads the gifter and giftee id from query params, and sends the gift if giftee is valid, returning the giftCard data.
func (h *ApiHandler) SendNewGift(response http.ResponseWriter, request *http.Request) {
	const spot = "api/SendNewGift"
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	gifterID, err := uuid.Parse(request.URL.Query().Get(GifterID))
	if err != nil {
		writeError(response, common.ExtendErrors(types.ErrGifterInvalid, err), 400)
		return
	}
	gifteeID, err := uuid.Parse(request.URL.Query().Get(GifteeID))
	if err != nil {
		writeError(response, common.ExtendErrors(types.ErrGifteeInvalid, err), 400)
		return
	}

	res, err := h.giftCardModule.IssueNewGiftCard(ctx, gifterID, gifteeID)
	if err != nil {
		writeError(response, err, 500)
		return
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		h.logger.Errorf("[%s] Failed to marshal json: %s", spot, err)
		writeError(response, err, 500)
		return
	}

	response.Write(jsonResp)
}

// RespondToGift responds to a received gift with either accepted or rejected
func (h *ApiHandler) RespondToGift(response http.ResponseWriter, request *http.Request) {
	const spot = "api/RespondToGift"
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	giftID, err := uuid.Parse(request.URL.Query().Get(GiftID))
	if err != nil {
		writeError(response, common.ExtendErrors(types.ErrGiftIDInvalid, err), 400)
		return
	}
	respStatus, err := parseResponse(request.URL.Query().Get(ResponseToGift))
	if err != nil {
		writeError(response, err, 400)
		return
	}

	err = h.giftCardModule.RespondToGift(ctx, giftID, respStatus)
	if err != nil {
		writeError(response, err, 500)
		return
	}

	response.WriteHeader(200)
}

// InquireSentGifts inquires gifts and if provided, will filter them based on status. Uses default pagination if none is provided.
func (h *ApiHandler) InquireSentGifts(response http.ResponseWriter, request *http.Request) {
	const spot = "api/InquireSentGifts"
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	gifterID, err := uuid.Parse(request.URL.Query().Get(GifterID))
	if err != nil {
		writeError(response, common.ExtendErrors(types.ErrGifterInvalid, err), 400)
		return
	}

	paginationData, err := parsePagination(request.URL.Query().Get(PageSize), request.URL.Query().Get(PageNumber))
	if err != nil {
		writeError(response, err, 400)
		return
	}

	wantedStatus, _ := parseStatus(request.URL.Query().Get(WantedStatus)) // we are fine with nil status

	res, err := h.giftCardModule.GetListOfSentGiftCards(ctx, gifterID, wantedStatus, *paginationData)
	if err != nil {
		writeError(response, err, 500)
		return
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		h.logger.Errorf("[%s] Failed to marshal json: %s", spot, err)
		writeError(response, err, 500)
		return
	}

	response.Write(jsonResp)
}

func (h *ApiHandler) InquireReceivedGifts(response http.ResponseWriter, request *http.Request) {
	const spot = "api/InquireReceivedGifts"
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	gifteeID, err := uuid.Parse(request.URL.Query().Get(GifteeID))
	if err != nil {
		writeError(response, common.ExtendErrors(types.ErrGifteeInvalid, err), 400)
		return
	}

	paginationData, err := parsePagination(request.URL.Query().Get(PageSize), request.URL.Query().Get(PageNumber))
	if err != nil {
		writeError(response, err, 400)
		return
	}

	wantedStatus, _ := parseStatus(request.URL.Query().Get(WantedStatus)) // we are fine with nil status

	res, err := h.giftCardModule.GetListOfReceivedGiftCards(ctx, gifteeID, wantedStatus, *paginationData)
	if err != nil {
		writeError(response, err, 500)
		return
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		h.logger.Errorf("[%s] Failed to marshal json: %s", spot, err)
		writeError(response, err, 500)
		return
	}

	response.Write(jsonResp)
}
