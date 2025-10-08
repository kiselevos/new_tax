package main

import (
	"flag"
	"log"
	"net"
	"new_tax/internal/server"

	pb "new_tax/gen/grpc/api"

	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	addr := flag.String("addr", ":50051", "listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("Problem with server", err)
	}

	srvGRPC := googleGRPC.NewServer()
	srv := server.NewGRPCServer()
	pb.RegisterTaxServiceServer(srvGRPC, srv)

	reflection.Register(srvGRPC)

	log.Printf("TaxService listening on %s", *addr)
	if err := srvGRPC.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
