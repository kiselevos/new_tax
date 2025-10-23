package client

import (
	"os"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func NewTaxClient() (pb.TaxServiceClient, *grpc.ClientConn, error) {
	addr := os.Getenv("BACKEND_ADDR")
	if addr == "" {
		if os.Getenv("IN_DOCKER") != "" {
			addr = "backend:50051"
		} else {
			addr = "localhost:50051"
		}
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return pb.NewTaxServiceClient(conn), conn, nil
}