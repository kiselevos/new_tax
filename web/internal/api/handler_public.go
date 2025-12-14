package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"github.com/kiselevos/new_tax/web/internal/metrics"
	"github.com/kiselevos/new_tax/web/internal/middleware"
)

type PublicHandler struct {
	TaxClient pb.TaxServiceClient
}

// NewPublicHandler - конструктор для PublicHandler
func NewPublicHandler(client pb.TaxServiceClient) *PublicHandler {
	return &PublicHandler{
		TaxClient: client,
	}
}

// HandlePublicCalc - публичный API
func (h *PublicHandler) HandlePublicCalc(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	log := logx.From(ctx)
	start := time.Now()

	region := middleware.GetRegion(ctx).Label
	client := "public"

	metrics.M.Calculator.Attempts.
		WithLabelValues(client, region).
		Inc()

	defer func() {
		metrics.M.Calculator.Duration.
			WithLabelValues(client, region).
			Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodPost {
		metrics.M.ErrorTypes.WithLabelValues(client, "method").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region).Inc()
		writeError(w, r, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PublicCalcRequest
	err := json.NewDecoder(r.Body).Decode(&dtoReq)
	if err != nil {
		log.Warn("api_invalid_json", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "json").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region).Inc()
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		log.Warn("validate_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "validate").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region).Inc()
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToProto()

	grpcResp, err := h.TaxClient.CalculatePublic(ctx, grpcReq)
	if err != nil {
		log.Warn("grpc_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "grpc").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region).Inc()
		writeError(w, r, "backend error", http.StatusInternalServerError)
		return
	}

	dtoResp := NewPublicResponseToJSON(grpcResp)

	metrics.M.Calculator.Success.WithLabelValues(client, region).Inc()

	writeJSON(w, http.StatusOK, dtoResp)
}
