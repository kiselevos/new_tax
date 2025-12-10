package api

import (
	"fmt"
	"html/template"
	"net/http"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

func RegisterPublicRoutes(mux *http.ServeMux, client pb.TaxServiceClient, apiVers string, tmpl *template.Template) {
	handlerPublic := NewPublicHandler(client)
	handlerPrivate := NewPrivateHandler(client)

	routePublic := fmt.Sprintf("/api/%s/calc", apiVers)
	routePrivate := fmt.Sprintf("/api/%s/private-calc", apiVers)

	mux.HandleFunc(routePublic, handlerPublic.HandlePublicCalc)
	mux.HandleFunc(routePrivate, handlerPrivate.HandlePrivateCalc)
}
