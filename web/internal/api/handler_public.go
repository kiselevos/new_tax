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

	region := middleware.GetRegion(ctx)
	client := "public"

	metrics.M.Calculator.Attempts.
		WithLabelValues(client, region.Label).
		Inc()

	defer func() {
		metrics.M.Calculator.Duration.
			WithLabelValues(client, region.Label).
			Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodPost {
		metrics.M.ErrorTypes.WithLabelValues(client, "method").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
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
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		log.Warn("validate_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "validate").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToProto()

	grpcResp, err := h.TaxClient.CalculatePublic(ctx, grpcReq)
	if err != nil {
		log.Warn("grpc_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "grpc").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, "backend error", http.StatusInternalServerError)
		return
	}

	dtoResp := NewPublicResponseToJSON(grpcResp)

	gross := float64(dtoReq.GrossSalary) / 100.0

	log.Info("business_calc",
		"client", client,
		"region", region.Name,
		"rid", middleware.GetRID(ctx),

		"gross_salary_rub", gross,
		"territorial_multiplier", dtoReq.TerritorialMultiplier,
		"northern_coefficient", dtoReq.NorthernCoefficient,
		"has_tax_privilege", false,
		"is_not_resident", false,

		"annual_tax", grpcResp.AnnualTaxAmount,
	)

	metrics.M.Calculator.GrossSalary.
		WithLabelValues(client, region.Label).
		Observe(gross)

	metrics.M.Calculator.Success.WithLabelValues(client, region.Label).Inc()
	writeJSON(w, http.StatusOK, dtoResp)
}
