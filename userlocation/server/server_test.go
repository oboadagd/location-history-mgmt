package grpcserver

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/oboadagd/location-history-mgmt/repository"
	"github.com/oboadagd/location-history-mgmt/service"
	"github.com/oboadagd/location-history-mgmt/testutils"
	"log"
	"net"

	pb "github.com/oboadagd/location-history-mgmt/userlocation/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var db *pg.DB

func init() {
	db := testutils.GetTestDB()

	testutils.CreateSchema(db)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := service.NewLocationService(locationRepository, locationHistoryRepository)

	pb.RegisterUserLocationServiceServer(s, &Server{
		LocationService: locationService,
		Context:         nil,
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
