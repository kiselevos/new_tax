package api

import (
	"encoding/json"
	"net/http"

	"github.com/kiselevos/new_tax/pkg/logx"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	_ = enc.Encode(v)

}

func writeError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	logx.From(r.Context()).Warn(
		"api_error",
		"status", status,
		"error", msg,
	)
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}
