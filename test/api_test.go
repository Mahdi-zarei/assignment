package test

import (
	"assignment/pkg/common"
	"assignment/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func createUsers(svc *ServiceTest) (user1 *types.UserData, user2 *types.UserData) {
	res1 := common.MustGetVal(http.Get("http://127.0.0.1:8080/register-user?username=Mahdi"))
	defer res1.Body.Close()

	res2 := common.MustGetVal(http.Get("http://127.0.0.1:8080/register-user?username=Ali"))
	defer res2.Body.Close()

	userData1 := types.UserData{}
	userData2 := types.UserData{}

	common.Must1(json.Unmarshal(common.MustGetVal(io.ReadAll(res1.Body)), &userData1))
	common.Must1(json.Unmarshal(common.MustGetVal(io.ReadAll(res2.Body)), &userData2))

	return &userData1, &userData2
}

func TestRegisterEndpoint(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()

	res, err := http.Get("http://127.0.0.1:8080/register-user?username=Mahdi")
	assert.Nil(t, err)
	defer res.Body.Close()

	userData := types.UserData{}
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(res.Body)), &userData)
	assert.Nil(t, err)

	resDB, err := svc.usersRepo.GetUserData(ctx, userData.ID)
	assert.Nil(t, err)
	assert.True(t, userData.Equals(*resDB))
}

func TestSendGiftEndpoint(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()

	userData1, userData2 := createUsers(svc)

	giftRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes.Body.Close()

	giftData := types.GiftCardData{}

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes.Body)), &giftData)
	assert.Nil(t, err)

	giftResDB, err := svc.giftCardRepo.GetGiftData(ctx, giftData.ID)
	assert.Nil(t, err)
	assert.True(t, giftData.Equals(*giftResDB))
}

func TestRespondToGiftEndpoint(t *testing.T) {
	svc := NewServiceTest()
	ctx := context.Background()

	userData1, userData2 := createUsers(svc)

	giftRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes.Body.Close()

	giftData := types.GiftCardData{}

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes.Body)), &giftData)
	assert.Nil(t, err)

	_, err = http.Get(fmt.Sprintf("http://127.0.0.1:8080/respond-to-gift?gift-id=%s&response-to-gift=accepted", giftData.ID))
	assert.Nil(t, err)

	giftResDB, err := svc.giftCardRepo.GetGiftData(ctx, giftData.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, giftResDB.Status, types.GiftCardStatusAccepted)
	assert.True(t, giftResDB.ResponseDate.After(time.Now().Add(-1*time.Minute)))
}

func TestInquireSentGiftsEndpoint(t *testing.T) {
	svc := NewServiceTest()

	userData1, userData2 := createUsers(svc)

	giftRes1, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes1.Body.Close()

	giftRes2, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes2.Body.Close()

	giftRes3, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes3.Body.Close()

	giftData1 := types.GiftCardData{}
	giftData2 := types.GiftCardData{}
	giftData3 := types.GiftCardData{}

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes1.Body)), &giftData1)
	assert.Nil(t, err)

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes2.Body)), &giftData2)
	assert.Nil(t, err)

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes3.Body)), &giftData3)
	assert.Nil(t, err)

	_, err = http.Get(fmt.Sprintf("http://127.0.0.1:8080/respond-to-gift?gift-id=%s&response-to-gift=accepted", giftData2.ID))
	assert.Nil(t, err)

	fullRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-sent-gifts?gifter-id=%s", userData1.ID))
	assert.Nil(t, err)
	defer fullRes.Body.Close()

	var fullResData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(fullRes.Body)), &fullResData)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(fullResData))

	resData1, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData1.ID
	})
	assert.True(t, found)
	resData2, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData2.ID
	})
	assert.True(t, found)
	resData3, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData3.ID
	})
	assert.True(t, found)

	assert.True(t, giftData1.Equals(*resData1))
	assert.True(t, giftData3.Equals(*resData3))
	assert.EqualValues(t, types.GiftCardStatusAccepted, resData2.Status)

	// test pagination
	pageRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-sent-gifts?gifter-id=%s&page-size=2&page-number=0", userData1.ID))
	assert.Nil(t, err)
	defer pageRes.Body.Close()

	var pageData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(pageRes.Body)), &pageData)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(pageData))

	// test status filter
	statRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-sent-gifts?gifter-id=%s&wanted-status=accepted", userData1.ID))
	assert.Nil(t, err)
	defer statRes.Body.Close()

	var statData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(statRes.Body)), &statData)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(statData))
}

func TestInquireReceivedGifts(t *testing.T) {
	svc := NewServiceTest()

	userData1, userData2 := createUsers(svc)

	giftRes1, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes1.Body.Close()

	giftRes2, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes2.Body.Close()

	giftRes3, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/send-new-gift?gifter-id=%s&giftee-id=%s", userData1.ID, userData2.ID))
	assert.Nil(t, err)
	defer giftRes3.Body.Close()

	giftData1 := types.GiftCardData{}
	giftData2 := types.GiftCardData{}
	giftData3 := types.GiftCardData{}

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes1.Body)), &giftData1)
	assert.Nil(t, err)

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes2.Body)), &giftData2)
	assert.Nil(t, err)

	err = json.Unmarshal(common.MustGetVal(io.ReadAll(giftRes3.Body)), &giftData3)
	assert.Nil(t, err)

	_, err = http.Get(fmt.Sprintf("http://127.0.0.1:8080/respond-to-gift?gift-id=%s&response-to-gift=accepted", giftData2.ID))
	assert.Nil(t, err)

	fullRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-received-gifts?giftee-id=%s", userData2.ID))
	assert.Nil(t, err)
	defer fullRes.Body.Close()

	var fullResData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(fullRes.Body)), &fullResData)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(fullResData))

	resData1, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData1.ID
	})
	assert.True(t, found)
	resData2, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData2.ID
	})
	assert.True(t, found)
	resData3, found := common.GetElem(fullResData, func(obj *types.GiftCardData) bool {
		return obj.ID == giftData3.ID
	})
	assert.True(t, found)

	assert.True(t, giftData1.Equals(*resData1))
	assert.True(t, giftData3.Equals(*resData3))
	assert.EqualValues(t, types.GiftCardStatusAccepted, resData2.Status)

	// test pagination
	pageRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-received-gifts?giftee-id=%s&page-size=2&page-number=0", userData2.ID))
	assert.Nil(t, err)
	defer pageRes.Body.Close()

	var pageData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(pageRes.Body)), &pageData)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(pageData))

	// test status filter
	statRes, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/inquire-received-gifts?giftee-id=%s&wanted-status=accepted", userData2.ID))
	assert.Nil(t, err)
	defer statRes.Body.Close()

	var statData []*types.GiftCardData
	err = json.Unmarshal(common.MustGetVal(io.ReadAll(statRes.Body)), &statData)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(statData))
}
