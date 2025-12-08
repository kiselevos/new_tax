package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/kiselevos/new_tax/gen/grpc/api"
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

	if r.Method != http.MethodPost {
		writeError(w, r, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dtoReq PrivateCalcRequest
	err := json.NewDecoder(r.Body).Decode(&dtoReq)
	if err != nil {
		writeError(w, r, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dtoReq.Validate(); err != nil {
		writeError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := dtoReq.ToPrivateProto()

	grpcResp, err := h.TaxClient.CalculatePrivate(ctx, grpcReq)
	if err != nil {
		writeError(w, r, "backend error", http.StatusInternalServerError)
		return
	}

	dtoResp := NewPrivateResponseToJSON(grpcResp)

	writeJSON(w, http.StatusOK, dtoResp)
}
