package api

import (
	"net/http"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

func RegisterPublicRoutes(mux *http.ServeMux, client pb.TaxServiceClient) {
	handler := NewPublicHandler(client)

	mux.HandleFunc("/api/v1/calc", handler.HandlePublicCalc)
}
