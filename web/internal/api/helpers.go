package api

import (
	"encoding/json"
	"net/http"

	"github.com/kiselevos/new_tax/pkg/logx"
	"google.golang.org/grpc/status"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	_ = enc.Encode(v)

}

func writeError(w http.ResponseWriter, r *http.Request, msg string, httpStatus int) {
	logx.From(r.Context()).Warn(
		"api_error",
		"status", httpStatus,
		"error", msg,
	)
	writeJSON(w, httpStatus, map[string]string{
		"error": msg,
	})
}

// grpcClientMsg возвращает сообщение об ошибке, безопасное для отдачи клиенту.
// 5xx — всегда generic, чтобы не утекали внутренние детали.
// 4xx — только gRPC-описание без префикса "rpc error: code = X desc = ...".
func grpcClientMsg(err error, httpStatus int) string {
	if httpStatus >= 500 {
		return "internal server error"
	}
	if st, ok := status.FromError(err); ok {
		return st.Message()
	}
	return "request error"
}
