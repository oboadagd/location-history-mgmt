// Package controller implements api layer for external clients.
// Through implementation of the LocationControllerInterface methods,
// it is possible to define validation and management of parameters
// it also invokes service layer.
package controller

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	respKit "github.com/oboadagd/kit-go/middleware/responses"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	"github.com/oboadagd/location-history-mgmt/service"
	"net/http"
	"time"
)

// LocationControllerInterface is the interface of Location controller layer. Contains definition of
// methods to manage the microservice apis.
type LocationControllerInterface interface {
	GetDistanceTraveled(c echo.Context) error
}

// LocationController represents the Location controller layer.
type LocationController struct {
	locationService service.LocationServiceInterface // Location service interface
}

// NewLocationController initializes Location controller layer.
func NewLocationController(locationService service.LocationServiceInterface) LocationControllerInterface {
	return &LocationController{
		locationService,
	}
}

// GetDistanceTraveled implements validation and management of parameters, then
// it invokes Location service layer of getting traveled distance by a username.
// Returns username data not found if username doesn't exist in Location model.
func (ctr *LocationController) GetDistanceTraveled(c echo.Context) error {
	var id, fd time.Time
	dateFormat := time.RFC3339
	un := c.Param("userName")

	log.Infof("REST Service GetDistanceTraveled started")

	if c.Param("initialDate") != "" {
		d, err := time.Parse(dateFormat, c.Param("initialDate"))
		if err != nil {
			return respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
		}
		id = d
	}

	if c.Param("finalDate") != "" {
		d, err := time.Parse(dateFormat, c.Param("finalDate"))
		if err != nil {
			return respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
		}
		fd = d
	}

	req := dto.GetDistanceTraveledRequest{
		UserName:    un,
		InitialDate: id,
		FinalDate:   fd,
	}

	vtr := validator.New()
	if err := vtr.RegisterValidation("patternazAZ09", dto.IsPatternUserName); err != nil {
		return respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	if err := vtr.RegisterValidation("maxDecimals", dto.IsMaxDecimals); err != nil {
		return respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	cvt := &dto.CustomValidatorSaveLoc{Validator: vtr}

	if err := cvt.Validate(req); err != nil {
		return respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	resp, err := ctr.locationService.GetDistanceTraveled(context.Background(), req)
	log.Infof("REST Service GetDistanceTraveled finished")
	if err != nil {
		log.Infof("err %v", err)
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
