// Package repository implements facade to relational database.
// Through implementation of the LocationHistoryRepositoryInterface
// methods, it is possible to define the necessary updates and fetches
// to manage LocationHistory entity model.
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

// LocationHistoryRepositoryInterface is the interface of LocationHistory repository layer.
// Contains definition of methods to manage the database representation of
// LocationHistory entity.
type LocationHistoryRepositoryInterface interface {
	Create(ctx context.Context, request dto.CreateLocationHistoryRequest) error
	GetDistanceByUserNameAndDateRange(ctx context.Context, request dto.GetDistanceTraveledRequest) (*dto.GetDistanceTraveledResponse, error)
	GetLastByUserName(ctx context.Context, request dto.GetLastByUserNameRequest) (*dto.GetLastByUserNameResponse, error)
}

// LocationHistoryRepository represents the relational database repository layer of
// LocationHistory entity. It's the historic registry of Location records.
type LocationHistoryRepository struct {
	db *pg.DB
}

// NewLocationHistoryRepository initializes repository of LocationHistory entity
func NewLocationHistoryRepository(db *pg.DB) LocationHistoryRepositoryInterface {
	return &LocationHistoryRepository{
		db,
	}
}

// Create implements insert action of LocationHistory entity
func (r *LocationHistoryRepository) Create(_ context.Context, request dto.CreateLocationHistoryRequest) error {

	lh := dto.LocationHistory{
		UserName:  request.UserName,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Distance:  request.Distance,
		UpdatedAt: time.Now(),
	}

	_, errIns := r.db.Model(&lh).Insert()
	if errIns != nil {
		return respKit.GenericBadRequestError(enums.ErrorInsertLocationCode, errIns.Error())
	}

	return nil
}

// GetDistanceByUserNameAndDateRange implements query select action of LocationHistory
// entity by username and date range. Returns the distance accumulated by a username across
// multiple records within a range of start date and end date. Returns error username data
// not found in case username doesn't exist.
func (r *LocationHistoryRepository) GetDistanceByUserNameAndDateRange(_ context.Context, request dto.GetDistanceTraveledRequest) (*dto.GetDistanceTraveledResponse, error) {
	var td []dto.GetDistanceTraveledResponse
	lh := dto.LocationHistory{}
	err := r.db.Model(&lh).
		Column("username").
		ColumnExpr("sum(distance) AS total_distance").
		Where("username = ?", request.UserName).
		Where("updated_at >= ?", request.InitialDate).
		Where("updated_at <= ?", request.FinalDate).
		Group("username").
		Select(&td)

	if err != nil {
		return &dto.GetDistanceTraveledResponse{}, respKit.GenericBadRequestError(enums.ErrorGetDistanceTraveledByUserNameCode, err.Error())
	}

	if len(td) == 0 {
		return &dto.GetDistanceTraveledResponse{}, respKit.GenericNotFoundError(enums.ErrorUserNameNotFoundCode, fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, request.UserName))
	}

	return &td[0], nil
}

// GetLastByUserName implements query select action of later LocationHistory
// entity by username. Returns later record by a username. Returns error username data
// not found in case username doesn't exist.
func (r *LocationHistoryRepository) GetLastByUserName(_ context.Context, request dto.GetLastByUserNameRequest) (*dto.GetLastByUserNameResponse, error) {
	var td []dto.GetLastByUserNameResponse
	lh := dto.LocationHistory{}
	err := r.db.Model(&lh).
		Column("username", "latitude", "longitude").
		Where("username = ?", request.UserName).
		Order("updated_at DESC").
		Select(&td)

	if err != nil {
		return &dto.GetLastByUserNameResponse{}, respKit.GenericBadRequestError(enums.ErrorGetLastLocationHistoryByUserNameCode, err.Error())
	}

	if len(td) == 0 {
		return &dto.GetLastByUserNameResponse{}, respKit.GenericNotFoundError(enums.ErrorUserNameNotFoundCode, fmt.Sprintf(enums.ErrorUserNameNotFoundMsg, request.UserName))
	}

	return &td[0], nil
}
