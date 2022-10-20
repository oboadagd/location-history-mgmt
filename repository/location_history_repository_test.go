package repository

import (
	"context"
	"fmt"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	"github.com/oboadagd/location-history-mgmt/testutils"
	"testing"
	"time"
)

func TestCreateLocationHistory(t *testing.T) {
	nameTest := "TestCreateLocationHistory"
	db = testutils.GetTestDB()
	defer db.Close()

	locationHistoryRepository := NewLocationHistoryRepository(db)
	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestCreateLocationHistory_ErrorInsertLocationCode(t *testing.T) {
	nameTest := "TestCreateLocationHistory_ErrorInsertLocationCode"
	db = testutils.GetTestDB()
	defer db.Close()

	locationHistoryRepository := NewLocationHistoryRepository(db)
	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh = &dto.CreateLocationHistoryRequest{}
	err = locationHistoryRepository.Create(ctx, *lh)

	if err == nil {
		t.Errorf("%s Expected %v but got %v", enums.ErrorInsertLocationCode, nameTest, err)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetDistanceByUserNameAndDateRange(t *testing.T) {
	nameTest := "TestGetDistanceByUserNameAndDateRange"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationHistoryRepository := NewLocationHistoryRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}
	end := time.Now()
	start := end.Add(-24 * time.Hour)
	dtr := dto.GetDistanceTraveledRequest{
		UserName:    lh.UserName,
		InitialDate: start,
		FinalDate:   end,
	}

	var resp *dto.GetDistanceTraveledResponse

	resp, err = locationHistoryRepository.GetDistanceByUserNameAndDateRange(ctx, dtr)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	if resp.TotalDistance == 0 {
		t.Errorf("%s: Expected %v but got %v", nameTest, 1, 0)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetDistanceByUserNameAndDateRange_ErrorUserNameNotFoundCode(t *testing.T) {
	nameTest := "TestGetDistanceByUserNameAndDateRange_ErrorUserNameNotFoundCode"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationHistoryRepository := NewLocationHistoryRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	dtr := dto.GetDistanceTraveledRequest{
		UserName: lh.UserName,
	}

	_, err = locationHistoryRepository.GetDistanceByUserNameAndDateRange(ctx, dtr)

	if err != nil && err.Error() != fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, lh.UserName) {
		t.Errorf("%s: Expected %v but got %v", nameTest, enums.ErrorUserNameNotFoundCode, err.Error())
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetLastByUserName(t *testing.T) {
	nameTest := "TestGetLastByUserName"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationHistoryRepository := NewLocationHistoryRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	llh := dto.GetLastByUserNameRequest{
		UserName: lh.UserName,
	}

	_, err = locationHistoryRepository.GetLastByUserName(ctx, llh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetLastByUserName_ErrorUserNameNotFoundCode(t *testing.T) {
	nameTest := "TestGetLastByUserName_ErrorUserNameNotFoundCode"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationHistoryRepository := NewLocationHistoryRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocationHistory()

	err = locationHistoryRepository.Create(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	llh := dto.GetLastByUserNameRequest{
		UserName: "usernamenotfound",
	}

	_, err = locationHistoryRepository.GetLastByUserName(ctx, llh)

	if err != nil && err.Error() != fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, llh.UserName) {
		t.Errorf("%s: Expected %v but got %v", nameTest, enums.ErrorUserNameNotFoundCode, err.Error())
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}
