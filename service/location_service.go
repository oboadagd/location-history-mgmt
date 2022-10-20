// Package service implements business logic layer.
// Through implementation of the LocationServiceInterface methods,
// it is possible to define business logic it also invokes repository
// layer.
package service

import (
	"context"
	geo "github.com/kellydunn/golang-geo"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-history-mgmt/repository"
	"time"
)

// LocationServiceInterface is the interface of Location service layer. Contains definition of
// methods to manage the business logic of Location and LocationHistory models.
type LocationServiceInterface interface {
	Save(ctx context.Context, request dto.SaveLocationRequest) error
	GetUsersByLocationAndRadius(ctx context.Context, request dto.GetUsersByLocationAndRadiusRequest) (*dto.GetUsersByLocationAndRadiusResponse, error)
	GetDistanceTraveled(ctx context.Context, request dto.GetDistanceTraveledRequest) (*dto.GetDistanceTraveledResponse, error)
}

// LocationService represents the Location service layer.
type LocationService struct {
	locationRepository        repository.LocationRepositoryInterface        // Location repository interface
	locationHistoryRepository repository.LocationHistoryRepositoryInterface // LocationHistory repository interface
}

// NewLocationService initializes Location service layer.
func NewLocationService(locationRepository repository.LocationRepositoryInterface, locationHistoryRepository repository.LocationHistoryRepositoryInterface) LocationServiceInterface {
	return &LocationService{
		locationRepository,
		locationHistoryRepository,
	}
}

// Save implements business logic of create and update actions of Location model.
// It creates records in Location and LocationHistory models when username doesn't
// already exist in Location model. Updates Location model and Create LocationHistory
// model record if username already exists in Location model. Sets traveled
// LocationHistory.distance from last to current location, if username already
// exists in Location model or zero otherwise.
func (s *LocationService) Save(ctx context.Context, request dto.SaveLocationRequest) error {

	var distance float64 = 0
	if s.locationRepository.ExistsByUserName(ctx, request.UserName) {
		if err := s.locationRepository.UpdateByUserName(ctx, request, request.UserName); err != nil {
			return err
		}

		llh := dto.GetLastByUserNameRequest{
			UserName: request.UserName,
		}

		resp, err := s.locationHistoryRepository.GetLastByUserName(ctx, llh)
		if err != nil {
			return err
		}
		ps := geo.NewPoint(resp.Latitude, resp.Longitude)
		pf := geo.NewPoint(request.Latitude, request.Longitude)
		distance = ps.GreatCircleDistance(pf)
	} else if err := s.locationRepository.Create(ctx, request); err != nil {
		return err
	}

	lh := dto.CreateLocationHistoryRequest{
		UserName:  request.UserName,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Distance:  distance,
	}

	if err := s.locationHistoryRepository.Create(ctx, lh); err != nil {
		return err
	}

	return nil
}

// GetUsersByLocationAndRadius implements business logic of getting a list of username's Location models
// that belongs to a given radius by requested page.
func (s *LocationService) GetUsersByLocationAndRadius(ctx context.Context, request dto.GetUsersByLocationAndRadiusRequest) (*dto.GetUsersByLocationAndRadiusResponse, error) {

	center := geo.NewPoint(request.Latitude, request.Longitude)
	pe := center.PointAtDistanceAndBearing(request.Radius, 90)
	pw := center.PointAtDistanceAndBearing(request.Radius, 270)
	pn := center.PointAtDistanceAndBearing(request.Radius, 0)
	ps := center.PointAtDistanceAndBearing(request.Radius, 180)

	llr := dto.GetByLatitudeLongitudeRangeRequest{
		LatitudeMin:  ps.Lat(),
		LatitudeMax:  pn.Lat(),
		LongitudeMin: pw.Lng(),
		LongitudeMax: pe.Lng(),
	}

	ulr, err := s.locationRepository.GetByLatitudeLongitudeRange(ctx, llr)
	if err != nil {
		return &dto.GetUsersByLocationAndRadiusResponse{}, err
	}

	var previewUsers []dto.Location
	for _, l := range ulr.Users {
		pf := geo.NewPoint(l.Latitude, l.Longitude)
		if center.GreatCircleDistance(pf) <= request.Radius {
			previewUsers = append(previewUsers, l)
		}
	}

	var resp dto.GetUsersByLocationAndRadiusResponse
	var minLimit = (request.Page - 1) * request.ItemsLimit
	var maxLimit = minLimit + request.ItemsLimit - 1
	var totalItems = uint64(len(ulr.Users))
	var totalPages = totalItems / request.ItemsLimit
	if totalItems%request.ItemsLimit != 0 {
		totalPages++
	}
	for i, l := range ulr.Users {
		if uint64(i) >= minLimit && uint64(i) <= maxLimit {
			resp.Users = append(resp.Users, l)
		}
	}

	resp.TotalItems = totalItems
	resp.TotalPages = totalPages

	return &resp, nil
}

// GetDistanceTraveled implements business logic of getting username traveled distance
// in a time range. Returns username data not found if username doesn't exist in Location model.
// If initial or final date has empty value then time range defaults to 1 day.
func (s *LocationService) GetDistanceTraveled(ctx context.Context, request dto.GetDistanceTraveledRequest) (*dto.GetDistanceTraveledResponse, error) {

	if request.FinalDate.IsZero() || request.InitialDate.IsZero() {
		end := time.Now()
		start := end.Add(-24 * time.Hour)
		request.InitialDate = start
		request.FinalDate = end
	}

	if request.FinalDate.Before(request.InitialDate) {
		request.InitialDate, request.FinalDate = request.FinalDate, request.InitialDate
	}

	dt, err := s.locationHistoryRepository.GetDistanceByUserNameAndDateRange(ctx, request)
	if err != nil {
		return &dto.GetDistanceTraveledResponse{}, err
	}

	return dt, nil
}
