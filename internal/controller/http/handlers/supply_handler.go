package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
)

type SupplyHandler struct {
	supplyService ports.SupplyService
}

func NewSupplyHandler(supplyService ports.SupplyService) *SupplyHandler {
	return &SupplyHandler{supplyService: supplyService}
}

func (h *SupplyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	supplies, err := h.supplyService.GetAll(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get supplies")
		return
	}

	response.Success(w, supplies)
}

func (h *SupplyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var supply domain.Supply
	if err := json.NewDecoder(r.Body).Decode(&supply); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if supply.IngredientID <= 0 {
		response.BadRequest(w, "ingredient_id must be > 0")
		return
	}
	if supply.Qty <= 0 {
		response.BadRequest(w, "qty must be > 0")
		return
	}
	if supply.SupplierName == "" {
		response.BadRequest(w, "supplier_name is required")
		return
	}

	if err := h.supplyService.Create(r.Context(), &supply); err != nil {
		response.InternalError(w, "failed to create supply")
		return
	}

	response.Created(w, supply)
}
