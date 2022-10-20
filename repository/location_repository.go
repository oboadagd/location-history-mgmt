// Package repository implements a facade to relational database.
// Through implementation of the LocationRepositoryInterface methods,
// it is possible to define the necessary updates and fetches to
// manage Location entity.
package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	respKit "github.com/oboadagd/kit-go/middleware/responses"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	"time"
)

// LocationRepositoryInterface is the interface of Location repository layer.
// Contains definition of methods to manage the database representation of
// Location entity.
type LocationRepositoryInterface interface {
	Create(ctx context.Context, request dto.SaveLocationRequest) error
	UpdateByUserName(ctx context.Context, request dto.SaveLocationRequest, userName string) error
	ExistsByUserName(ctx context.Context, userName string) bool
	GetByLatitudeLongitudeRange(ctx context.Context, request dto.GetByLatitudeLongitudeRangeRequest) (*dto.GetUsersByLocationAndRadiusResponse, error)
}

// LocationRepository  represents the relational database repository layer of
// Location entity. It's the registry of Location entity records. Exists a
// unique record for each username. Username's Location registry is updated
// each time geographic coordinates change.
type LocationRepository struct {
	Db *pg.DB // available database
}

// NewLocationRepository initializes repository of Location entity.
func NewLocationRepository(db *pg.DB) LocationRepositoryInterface {
	return &LocationRepository{
		db,
	}
}

// Create implements insert action of Location entity.
func (r *LocationRepository) Create(_ context.Context, request dto.SaveLocationRequest) error {

	l := dto.Location{
		UserName:  request.UserName,
		Latitude:  request.Latitude,
		Longitude: request.Latitude,
		UpdatedAt: time.Now(),
	}
	_, errIns := r.Db.Model(&l).Insert()
	if errIns != nil {
		return respKit.GenericBadRequestError(enums.ErrorInsertLocationCode, errIns.Error())
	}

	return nil
}

// UpdateByUserName implements update action of Location entity by username.
// Returns username data not found if username doesn't exist.
func (r *LocationRepository) UpdateByUserName(_ context.Context, request dto.SaveLocationRequest, userName string) error {
	var resp []dto.Location
	err := r.Db.Model(&dto.Location{}).Where("userName = ?", userName).Select(&resp)

	if resp == nil || (err != nil && err == pg.ErrNoRows) {
		return respKit.GenericNotFoundError(enums.ErrorUserNameNotFoundCode, fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, userName))
	}

	resp[0].Latitude = request.Latitude
	resp[0].Longitude = request.Longitude
	resp[0].UpdatedAt = time.Now()

	if _, err := r.Db.Model(&resp[0]).Where("userName = ?", userName).Update(); err != nil {
		return respKit.GenericBadRequestError(enums.ErrorUpdateLocationCode, err.Error())
	}

	return nil
}

// ExistsByUserName implements exist action of Location entity by username.
// Returns true if successful and false otherwise.
func (r *LocationRepository) ExistsByUserName(_ context.Context, userName string) bool {
	var resp []dto.Location
	err := r.Db.Model(&dto.Location{}).Where("userName = ?", userName).Select(&resp)
	if resp == nil || (err != nil && err == pg.ErrNoRows) {
		return false
	}

	return true
}

// GetByLatitudeLongitudeRange implements query select action of Location entity on
// a square area. The square area is defined by the maximum and minimum latitude and
// also by the maximum and minimum longitude.
func (r *LocationRepository) GetByLatitudeLongitudeRange(_ context.Context, request dto.GetByLatitudeLongitudeRangeRequest) (*dto.GetUsersByLocationAndRadiusResponse, error) {
	var u []dto.Location
	lr := dto.GetUsersByLocationAndRadiusResponse{}
	l := dto.Location{}
	err := r.Db.Model(&l).
		Where("latitude >= ?", request.LatitudeMin).
		Where("latitude <= ?", request.LatitudeMax).
		Where("longitude >= ?", request.LongitudeMin).
		Where("longitude <= ?", request.LongitudeMax).
		Select(&u)

	lr.Users = u

	if err != nil {
		return &lr, respKit.GenericBadRequestError(enums.ErrorGetByLatitudeLongitudeRangeMsg, err.Error())
	}

	return &lr, nil
}
