package client

import (
	pb "github.com/kiselevos/new_tax/gen/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTaxClient(addr string) (pb.TaxServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return pb.NewTaxServiceClient(conn), conn, nil
}
