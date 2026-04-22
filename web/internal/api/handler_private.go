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

	ctx := r.Context()
	log := logx.From(ctx)

	start := time.Now()
	region := middleware.GetRegion(ctx)
	client := "private"

	metrics.M.Calculator.Attempts.
		WithLabelValues(client, region.Label).
		Inc()
	defer func() {
		metrics.M.Calculator.Duration.
			WithLabelValues(client, region.Label).
			Observe(time.Since(start).Seconds())
	}()

	// Validate method
	if r.Method != http.MethodPost {
		metrics.M.ErrorTypes.WithLabelValues(client, "method").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("x-api-key")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PrivateCalcRequest
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		log.Warn("api_invalid_json", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "json").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		log.Warn("api_validation_failed", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "validation").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToPrivateProto()
	md := metadata.Pairs("x-api-key", apiKey)
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcResp, err := h.TaxClient.CalculatePrivate(ctx, grpcReq)
	if err != nil {
		log.Warn("grpc_error", "err", err)
		metrics.M.ErrorTypes.WithLabelValues(client, "grpc").Inc()
		metrics.M.Calculator.Failed.WithLabelValues(client, region.Label).Inc()
		httpStatus := grpcToHTTP(err)
		writeError(w, r, grpcClientMsg(err, httpStatus), httpStatus)
		return
	}

	dtoResp := NewPrivateResponseToJSON(grpcResp)

	gross := float64(dtoReq.GrossSalary) / 100.0

	log.Info("business_calc",
		"client", client,
		"region", region.Name,
		"rid", middleware.GetRID(ctx),

		"gross_salary_rub", gross,
		"employment_type", dtoReq.EmploymentType,
		"territorial_multiplier", dtoReq.TerritorialMultiplier,
		"northern_coefficient", dtoReq.NorthernCoefficient,
		"has_tax_privilege", dtoReq.HasTaxPrivilege,
		"is_not_resident", dtoReq.IsNotResident,

		"annual_tax", grpcResp.AnnualTaxAmount,
	)

	metrics.M.Calculator.GrossSalary.
		WithLabelValues(client, region.Label).
		Observe(gross)

	metrics.M.Calculator.Success.WithLabelValues(client, region.Label).Inc()
	writeJSON(w, http.StatusOK, dtoResp)
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
