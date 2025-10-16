package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
	"github.com/go-chi/chi/v5"
)

type TableHandler struct {
	tableService ports.TableService
}

func NewTableHandler(tableService ports.TableService) *TableHandler {
	return &TableHandler{tableService: tableService}
}

func (h *TableHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tables, err := h.tableService.GetAll(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get tables")
		return
	}

	response.Success(w, tables)
}

func (h *TableHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid table id")
		return
	}

	table, err := h.tableService.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w, "failed to get table")
		return
	}

	if table == nil {
		response.NotFound(w, "table not found")
		return
	}

	response.Success(w, table)
}

func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var table domain.Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.tableService.Create(r.Context(), &table); err != nil {
		response.InternalError(w, "failed to create table")
		return
	}

	response.Created(w, table)
}

func (h *TableHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid table id")
		return
	}

	var req struct {
		Status domain.TableStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.tableService.UpdateStatus(r.Context(), id, req.Status); err != nil {
		response.InternalError(w, "failed to update table status")
		return
	}

	response.Success(w, map[string]string{"message": "table status updated"})
}

func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid table id")
		return
	}

	if err := h.tableService.Delete(r.Context(), id); err != nil {
		response.InternalError(w, "failed to delete table")
		return
	}

	response.Success(w, map[string]string{"message": "table deleted"})
}
