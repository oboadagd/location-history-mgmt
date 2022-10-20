// Package grpcserver implements routines to start up a grpc server.
// After the server is up and running can listen to request sent by grpc clients.
// It also implements routines that invoke service layer and return a http response
// to callers(grpc clients).
package grpcserver

import (
	"github.com/labstack/echo/v4"
	"github.com/oboadagd/location-history-mgmt/service"
	pb "github.com/oboadagd/location-history-mgmt/userlocation/proto"
)

// Server is the representation of a grpc server state.
type Server struct {
	pb.UserLocationServiceServer                                  // interface is the server API for UserLocationService
	LocationService              service.LocationServiceInterface // interface is the local service layer
	Context                      echo.Context                     // the context of the current HTTP request
}
