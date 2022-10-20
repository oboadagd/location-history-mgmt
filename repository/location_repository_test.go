package repository

import (
	"context"
	"fmt"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	"testing"
)
import "github.com/oboadagd/location-history-mgmt/testutils"
import "github.com/go-pg/pg/v10"

var db *pg.DB

func TestCreateLocation(t *testing.T) {
	nameTest := "TestCreateLocation"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := NewLocationRepository(db)
	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

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

func TestCreateLocation_ErrorInsertLocationCode(t *testing.T) {
	nameTest := "TestCreateLocation_ErrorInsertLocationCode"
	db = testutils.GetTestDB()
	defer db.Close()

	locationRepository := NewLocationRepository(db)
	ctx := context.Background()

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	err = locationRepository.Create(ctx, *l)

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

func TestExistsByUserName_True(t *testing.T) {
	nameTest := "TestExistsByUserName_True"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationRepository := NewLocationRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	var resp bool

	resp = locationRepository.ExistsByUserName(ctx, l.UserName)

	if !resp {
		t.Errorf("%s: Expected %v but got %v", nameTest, true, resp)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestExistsByUserName_False(t *testing.T) {
	nameTest := "TestExistsByUserName_False"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationRepository := NewLocationRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	var resp bool
	resp = locationRepository.ExistsByUserName(ctx, "usernamenotfound")

	if resp {
		t.Errorf("%s: Expected %v but got %v", nameTest, false, resp)
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}

func TestUpdateByUserName_ErrorUserNameNotFoundCode(t *testing.T) {
	nameTest := "TestUpdateByUserName_ErrorUserNameNotFoundCode"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationRepository := NewLocationRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	err = locationRepository.UpdateByUserName(ctx, *l, "usernamenotfound")
	if err != nil && err.Error() != fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, "usernamenotfound") {
		t.Errorf("%s: Expected %v but got %v", nameTest, enums.ErrorUserNameNotFoundCode, err.Error())
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%v: %v", nameTest, err)
		return
	}

	t.Logf("%v Success", nameTest)
}

func TestUpdateByUserName_ErrorUpdateLocationCode(t *testing.T) {
	nameTest := "TestUpdateByUserName_ErrorUpdateLocationCode"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationRepository := NewLocationRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	userName := l.UserName
	l = &dto.SaveLocationRequest{}
	err = locationRepository.UpdateByUserName(ctx, *l, userName)
	if err != nil && err.Error() == fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, userName) {
		t.Errorf("%s: Expected %v but got %v", nameTest, enums.ErrorUpdateLocationCode, err.Error())
		return
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%v: %v", nameTest, err)
		return
	}

	t.Logf("%v Success", nameTest)
}

func TestGetByLatitudeLongitudeRange(t *testing.T) {
	nameTest := "TestGetByLatitudeLongitudeRange"
	db = testutils.GetTestDB()
	defer db.Close()

	ctx := context.Background()
	locationRepository := NewLocationRepository(db)

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	l := testutils.GetLocation()

	err = locationRepository.Create(ctx, *l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	llr := dto.GetByLatitudeLongitudeRangeRequest{
		LongitudeMax: 20,
		LongitudeMin: 1,
		LatitudeMax:  20,
		LatitudeMin:  1,
	}

	var resp *dto.GetUsersByLocationAndRadiusResponse

	resp, err = locationRepository.GetByLatitudeLongitudeRange(ctx, llr)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	if len(resp.Users) == 0 {
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
