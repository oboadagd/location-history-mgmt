package grpcserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/oboadagd/location-history-mgmt/service"
	"net"

	pb "github.com/oboadagd/location-history-mgmt/userlocation/proto"

	"google.golang.org/grpc"
)

var host string = "0.0.0.0"         // host ip of the local grpc server
var port string = "50061"           // host port of the local grpc server
var addr string = host + ":" + port // host base url of the local grpc server

// GrpcServe starts up the grpc server.
func GrpcServe(locationService service.LocationServiceInterface, ctx echo.Context) error {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	defer lis.Close()
	log.Infof("grpc server on %s", lis.Addr().String())

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	pb.RegisterUserLocationServiceServer(s, &Server{
		LocationService: locationService,
		Context:         ctx,
	})

	defer s.Stop()
	return s.Serve(lis)
}
