package grpcserver

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	pb "github.com/oboadagd/location-history-mgmt/userlocation/proto"
)

func (s *Server) SaveLocation(ctx context.Context, req *pb.SaveLocationRequest) (*pb.SaveLocationResponse, error) {

	log.Infof("GRPC SaveLocation started: %v", req)

	inReq := dto.SaveLocationRequest{
		UserName:  req.UserName,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := s.LocationService.Save(ctx, inReq); err != nil {
		log.Fatalf("SaveLocation error: %+v ", err)
		return &pb.SaveLocationResponse{}, err
	}

	log.Infof("GRPC SaveLocation finished: ")
	return &pb.SaveLocationResponse{Message: enums.LocationCreated}, nil
}

func (s *Server) GetUsersByLocationAndRadius(ctx context.Context, req *pb.GetUsersByLocationAndRadiusRequest) (*pb.GetUsersByLocationAndRadiusResponse, error) {

	log.Infof("GRPC GetUsersByLocationAndRadius started: %v", req)

	var pbResp = pb.GetUsersByLocationAndRadiusResponse{}

	inReq := dto.GetUsersByLocationAndRadiusRequest{
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Radius:     req.Radius,
		Page:       req.Page,
		ItemsLimit: req.ItemsLimit,
	}

	resp, err := s.LocationService.GetUsersByLocationAndRadius(ctx, inReq)

	if err != nil {
		log.Fatalf("GRPC GetUsersByLocationAndRadius error, %+v ", err)
		return &pb.GetUsersByLocationAndRadiusResponse{}, err
	}

	for _, u := range resp.Users {
		pbResp.Users = append(pbResp.Users, &pb.Location{
			UserName:  u.UserName,
			Latitude:  u.Latitude,
			Longitude: u.Longitude,
		})
	}

	pbResp.TotalPages = resp.TotalPages
	pbResp.TotalItems = resp.TotalItems

	log.Infof("GRPC GetUsersByLocationAndRadius finished: ")
	return &pbResp, nil
}
