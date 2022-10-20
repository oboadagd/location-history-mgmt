package controller

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/oboadagd/location-history-mgmt/repository"
	"github.com/oboadagd/location-history-mgmt/service"
	"github.com/oboadagd/location-history-mgmt/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var db *pg.DB

func TestGetDistanceTraveled(t *testing.T) {
	nameTest := "GetDistanceTraveled"
	db = testutils.GetTestDB()
	defer db.Close()

	ctxBkg := context.Background()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := service.NewLocationService(locationRepository, locationHistoryRepository)
	locationController := NewLocationController(locationService)

	dateFormat := "%d-%02d-%02dT%02d:%02d:%02d+00:00"
	base := time.Now()
	end := base.Add(24 * time.Hour)
	start := base.Add(-24 * time.Hour)

	startStr := fmt.Sprintf(dateFormat,
		start.Year(), start.Month(), start.Day(),
		start.Hour(), start.Minute(), start.Second())

	endStr := fmt.Sprintf(dateFormat,
		end.Year(), end.Month(), end.Day(),
		end.Hour(), end.Minute(), end.Second())

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/location-history-mgmt/locations/distance/:userName/:initialDate/:finalDate")

	err := testutils.CreateSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh := testutils.GetLocation()

	err = locationService.Save(ctxBkg, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	lh = testutils.GetLocation()
	lh.Latitude = 40
	lh.Longitude = 40
	err = locationService.Save(ctxBkg, *lh)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	type test struct {
		data           []string
		resultValidate []string
		answer         string
	}

	tests := []test{
		{[]string{"usernamesample", startStr, endStr}, []string{""}, "success"},
		{[]string{"usernamesample", "", ""}, []string{""}, "success"},
		{[]string{"", startStr, endStr}, []string{"username", "required"}, "userName required failed"},
		{[]string{"usr", startStr, endStr}, []string{"username", "min"}, "userName min failed"},
		{[]string{"username123456789", startStr, endStr}, []string{"username", "max"}, "userName max failed"},
		{[]string{"username_1", startStr, endStr}, []string{"username", "pattern"}, "userName pattern failed"},
		{[]string{"usernamesample", "10-10-2022", endStr}, []string{"time"}, "date failed"},
		{[]string{"usernamesample", startStr, "10-10-2022"}, []string{"time"}, "date failed"},
	}

	for _, v := range tests {
		ctx.SetParamNames("userName", "initialDate", "finalDate")
		ctx.SetParamValues(v.data[0], v.data[1], v.data[2])

		err = locationController.GetDistanceTraveled(ctx)

		if err != nil && testutils.EvaluateErrConditions(err.Error(), v.resultValidate) {
			t.Errorf("%s: Expected %v but got %v", nameTest, v.answer, err.Error())
			return
		}
	}

	err = testutils.DropSchema(db)
	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	t.Logf("%s Success", nameTest)
}
