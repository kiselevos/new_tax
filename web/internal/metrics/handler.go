package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler() http.Handler {
	return promhttp.Handler()
}

func Mount(mux *http.ServeMux) {
	mux.Handle("/metrics", Handler())
}
