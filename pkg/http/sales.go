package http

import (
	"context"
	"dataflow-api/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type SaleService interface {
	GetAll(context.Context) ([]*model.Sale, error)
	Create(context.Context, *model.Sale) error
	Calculate(context.Context, string, *time.Time, *time.Time) (float64, error)
}

type Handler struct {
	Service SaleService
}

func NewSaleHandler(r *mux.Router, s SaleService) {
	handler := &Handler{
		Service: s,
	}
	r.HandleFunc("/data", handler.CreateSale).Methods("POST")
	r.HandleFunc("/data", handler.GetAllSales).Methods("GET")
	r.HandleFunc("/calculate", handler.Calculate).Methods("POST")
}

func (h *Handler) CreateSale(w http.ResponseWriter, r *http.Request) {
	var s model.Sale

	if err := decodeJSONBody(w, r, &s); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{Status: "error", Error: err.Error()})
		return
	}

	if err := model.ValidateSale(&s); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{Status: "error", Error: err.Error()})
		return
	}

	if err := h.Service.Create(r.Context(), &s); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{Status: "error", Error: "Failed to create sale: " + err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, successResponse{Status: "success"})
}

func (h *Handler) GetAllSales(w http.ResponseWriter, r *http.Request) {
	sales, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{Status: "error", Error: "Failed to fetch sales: " + err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, sales)
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest

	if err := decodeJSONBody(w, r, &req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{Status: "error", Error: err.Error()})
		return
	}

	if req.Operation != operationTotalSales {
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{Status: "error", Error: "invalid operation"})
		return
	}

	sum, err := h.Service.Calculate(r.Context(), req.StoreID, req.StartDate, req.EndDate)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{Status: "error", Error: "Calculation error: " + err.Error()})
		return
	}

	response := CalculateResponse{
		StoreID:    req.StoreID,
		TotalSales: sum,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}

	writeJSONResponse(w, http.StatusOK, response)
}
