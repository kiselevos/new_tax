package api

import (
	"fmt"
	"net/http"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

func RegisterPublicRoutes(mux *http.ServeMux, client pb.TaxServiceClient, apiVers string) {
	handler := NewPublicHandler(client)
	route := fmt.Sprintf("/api/%s/calc", apiVers)

	mux.HandleFunc(route, handler.HandlePublicCalc)
}
