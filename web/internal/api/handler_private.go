package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web/internal/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PrivateHandler struct {
	TaxClient pb.TaxServiceClient
}

// NewPrivateHandler - конструктор для PrivateHandler
func NewPrivateHandler(client pb.TaxServiceClient) *PrivateHandler {
	return &PrivateHandler{
		TaxClient: client,
	}
}

func (h *PrivateHandler) HandlePrivateCalc(w http.ResponseWriter, r *http.Request) {

	log := logx.From(r.Context())
	start := time.Now()

	metrics.M.PrivateAPI.Attempts.Inc()
	defer func() {
		metrics.M.PrivateAPI.Duration.Observe(time.Since(start).Seconds())
	}()

	// Validate method
	if r.Method != http.MethodPost {
		metrics.M.ErrorTypes.WithLabelValues("private", "method").Inc()
		metrics.M.PrivateAPI.Failed.Inc()
		writeError(w, r, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("x-api-key")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PrivateCalcRequest
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		log.Warn("api_invalid_json", "err", err)
		metrics.M.ErrorTypes.WithLabelValues("private", "json").Inc()
		metrics.M.PrivateAPI.Failed.Inc()
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		log.Warn("api_validation_failed", "err", err)
		metrics.M.ErrorTypes.WithLabelValues("private", "validation").Inc()
		metrics.M.PrivateAPI.Failed.Inc()
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToPrivateProto()
	md := metadata.Pairs("x-api-key", apiKey)
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcResp, err := h.TaxClient.CalculatePrivate(ctx, grpcReq)
	if err != nil {
		log.Warn("grpc_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues("private", "grpc").Inc()
		metrics.M.PrivateAPI.Failed.Inc()
		writeError(w, r, err.Error(), grpcToHTTP(err))
		return
	}

	dtoResp := NewPrivateResponseToJSON(grpcResp)

	writeJSON(w, http.StatusOK, dtoResp)

	metrics.M.PrivateAPI.Success.Inc()
}

func grpcToHTTP(err error) int {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
