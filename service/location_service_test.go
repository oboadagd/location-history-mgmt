package service

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-history-mgmt/repository"
	"github.com/oboadagd/location-history-mgmt/testutils"
	"testing"
	"time"
)

var db *pg.DB

func TestSave_Create(t *testing.T) {
	nameTest := "TestSave_Create"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := NewLocationService(locationRepository, locationHistoryRepository)

	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationService.Save(ctx, *l)

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

func TestSave_Update(t *testing.T) {
	nameTest := "TestSave_Update"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := NewLocationService(locationRepository, locationHistoryRepository)

	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocation()

	err = locationService.Save(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh.Longitude = 20
	err = locationService.Save(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	llh := dto.GetLastByUserNameRequest{
		UserName: lh.UserName,
	}

	resp, err := locationHistoryRepository.GetLastByUserName(ctx, llh)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	if resp.Longitude != 20 {
		t.Errorf("%s: Expected %v but got %v", nameTest, 20, resp.Longitude)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetUsersByLocationAndRadius(t *testing.T) {
	nameTest := "TestGetUsersByLocationAndRadius"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := NewLocationService(locationRepository, locationHistoryRepository)

	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocation()

	err = locationService.Save(ctx, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	ulr := dto.GetUsersByLocationAndRadiusRequest{
		Latitude:   lh.Latitude,
		Longitude:  lh.Longitude,
		Radius:     10,
		Page:       1,
		ItemsLimit: 10,
	}

	resp, err := locationService.GetUsersByLocationAndRadius(ctx, ulr)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	if len(resp.Users) != 1 {
		t.Errorf("%s: Expected %v but got %v", nameTest, 1, len(resp.Users))
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestGetDistanceTraveled(t *testing.T) {
	nameTest := "GetDistanceTraveled"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := NewLocationService(locationRepository, locationHistoryRepository)

	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocation()

	err = locationService.Save(ctx, *lh)

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

	resp, err := locationService.GetDistanceTraveled(ctx, dtr)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	if resp.TotalDistance != 0 {
		t.Errorf("%s: Expected %v but got %v", nameTest, 0, resp.TotalDistance)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}
