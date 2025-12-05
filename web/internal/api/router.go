package api

import (
	"fmt"
	"html/template"
	"net/http"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
)

func RegisterPublicRoutes(mux *http.ServeMux, client pb.TaxServiceClient, apiVers string, tmpl *template.Template) {
	handlerPublic := NewPublicHandler(client)
	handlerDocs := NewDocsHandler(apiVers, tmpl)

	route := fmt.Sprintf("/api/%s/calc", apiVers)

	handlerDocs.AddEndpoint(EndpointInfo{
		Method:             "POST",
		Path:               route,
		Description:        "Открытый базовый расчёт НДФЛ",
		ExampleRequestLite: `{"gross_salary": 120000}`,
		ExampleRequestFull: `{
  "gross_salary": 120000,
  "territorial_multiplier": 120,
  "northern_coefficient": 150
}`,
		ExampleResponse: `{
  "gross_salary": 120000,
  "territorial_multiplier": 120,
  "northern_coefficient": 150,

  "annual_tax_amount": 17160,
  "annual_gross_income": 180000,
  "annual_net_income": 162840,

  "monthly_details": [
    {
      "month": "2024-01-01T00:00:00Z",
      "monthly_gross_income": 180000,
      "monthly_net_income": 162840,
      "monthly_tax_amount": 17160,

      "annual_gross_income": 180000,
      "annual_net_income": 162840,
      "annual_tax_amount": 17160
    }
  ]
}`,
		Group: "public",
	})
	mux.HandleFunc("/api/docs", handlerDocs.HandleApiDocs)

	mux.HandleFunc(route, handlerPublic.HandlePublicCalc)
}
