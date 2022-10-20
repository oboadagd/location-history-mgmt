//go:build integration
// +build integration

package grpcserver

import (
	"context"
	"github.com/oboadagd/location-common/enums"
	"testing"

	pb "github.com/oboadagd/location-history-mgmt/userlocation/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSaveLocation(t *testing.T) {
	nameTest := "TestSaveLocation"
	ctx := context.Background()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), creds)

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserLocationServiceClient(conn)

	type testReq struct {
		dataStr  string
		dataNumb []float64
		answer   string
	}

	tests := []testReq{
		{"usernamesample", []float64{10, 10}, enums.LocationCreated},
		{"usernamesample", []float64{20, 20}, enums.LocationCreated},
	}

	for _, v := range tests {

		req := &pb.SaveLocationRequest{
			UserName:  v.dataStr,
			Latitude:  v.dataNumb[0],
			Longitude: v.dataNumb[1],
		}
		resp, err := c.SaveLocation(ctx, req)

		if err != nil {
			t.Errorf("%s: unexpected error %v", nameTest, err)
			return
		}

		if resp.Message != enums.LocationCreated {
			t.Errorf("%s: Expected %v but got %v", nameTest, enums.LocationCreated, resp.Message)
			return
		}
	}
}

func TestGetUsersByLocationAndRadius(t *testing.T) {
	nameTest := "TestSaveLocation"
	ctx := context.Background()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), creds)

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserLocationServiceClient(conn)

	l := &pb.SaveLocationRequest{
		UserName:  "username",
		Latitude:  10,
		Longitude: 10,
	}
	_, err = c.SaveLocation(ctx, l)

	if err != nil {
		t.Errorf("%s: %v", nameTest, err)
		return
	}

	type testReq struct {
		data    []float64
		dataInt []uint64
		answer  int
	}

	tests := []testReq{
		{[]float64{10, 10, 10}, []uint64{1, 10}, 1},
	}

	for _, v := range tests {
		req := &pb.GetUsersByLocationAndRadiusRequest{
			Latitude:   v.data[0],
			Longitude:  v.data[1],
			Radius:     v.data[2],
			Page:       v.dataInt[0],
			ItemsLimit: v.dataInt[1],
		}

		resp, err := c.GetUsersByLocationAndRadius(ctx, req)

		if err != nil {
			t.Errorf("%s: %v", nameTest, err)
			return
		}

		if len(resp.Users) != v.answer {
			t.Errorf("%s: Expected %v but got %v", nameTest, v.answer, len(resp.Users))
			return
		}
	}
}
