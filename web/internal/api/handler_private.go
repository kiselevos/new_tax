package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
	"github.com/kiselevos/new_tax/pkg/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PrivateHandler struct {
	TaxClient pb.TaxServiceClient
}

// NewPublicHandler - конструктор для PublicHandler
func NewPrivateHandler(client pb.TaxServiceClient) *PrivateHandler {
	return &PrivateHandler{
		TaxClient: client,
	}
}

// HandlePrivateCalc - приватный API
func (h *PrivateHandler) HandlePrivateCalc(w http.ResponseWriter, r *http.Request) {

	log := logx.From(r.Context())

	if r.Method != http.MethodPost {
		writeError(w, r, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("x-api-key")

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PrivateCalcRequest
	err := json.NewDecoder(r.Body).Decode(&dtoReq)
	if err != nil {
		log.Warn("api_invalid_json", "err", err)
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		log.Warn("api_validation_failed", "err", err)
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToPrivateProto()

	md := metadata.Pairs("x-api-key", apiKey)
	ctx = metadata.NewOutgoingContext(ctx, md)

	grpcResp, err := h.TaxClient.CalculatePrivate(ctx, grpcReq)
	if err != nil {
		writeError(w, r, err.Error(), grpcToHTTP(err))
		return
	}

	dtoResp := NewPrivateResponseToJSON(grpcResp)

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
