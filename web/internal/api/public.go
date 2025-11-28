package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
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

	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PublicCalcRequest
	err := json.NewDecoder(r.Body).Decode(&dtoReq)
	if err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToProto()

	grpcResp, err := h.TaxClient.CalculatePublic(ctx, grpcReq)
	if err != nil {
		writeError(w, "backend error", http.StatusInternalServerError)
		return
	}

	dtoResp := NewPublicResponseToJSON(grpcResp)

	writeJSON(w, http.StatusOK, dtoResp)
}
