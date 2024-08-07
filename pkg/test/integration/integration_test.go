package integration

import (
	"bytes"
	servicehttp "dataflow-api/pkg/http"
	"dataflow-api/pkg/model"
	repositorysale "dataflow-api/pkg/repository/sale"
	"dataflow-api/pkg/service/sale"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouter() *mux.Router {
	repo := repositorysale.NewInMemorySaleRepository()
	svc := sale.NewService(repo)
	handler := servicehttp.Handler{Service: svc}

	router := mux.NewRouter()
	router.HandleFunc("/data", handler.CreateSale).Methods("POST")
	router.HandleFunc("/data", handler.GetAllSales).Methods("GET")
	router.HandleFunc("/calculate", handler.Calculate).Methods("POST")

	return router
}

func TestIntegration_Calculate_WithDateRange(t *testing.T) {
	router := setupRouter()

	now := time.Now()
	sale1 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     now,
	}
	sale2 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     now.Add(-2 * time.Hour),
	}

	createSale(router, sale1, t)
	createSale(router, sale2, t)

	startDate := now.Add(-1 * time.Hour)
	endDate := now.Add(1 * time.Hour)

	reqBody := servicehttp.CalculateRequest{
		Operation: "total_sales",
		StoreID:   "store1",
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res servicehttp.CalculateResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res.TotalSales != 100.0 {
		t.Fatalf("expected total sales to be 100.0, got %f", res.TotalSales)
	}
}

func TestIntegration_Calculate_WithoutDates(t *testing.T) {
	router := setupRouter()

	now := time.Now()
	sale := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     now,
	}

	createSale(router, sale, t)

	reqBody := servicehttp.CalculateRequest{
		Operation: "total_sales",
		StoreID:   "store1",
		StartDate: nil,
		EndDate:   nil,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res servicehttp.CalculateResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res.TotalSales != 100.0 {
		t.Fatalf("expected total sales to be 100.0, got %f", res.TotalSales)
	}
}

func TestIntegration_Calculate_NoSales(t *testing.T) {
	router := setupRouter()

	reqBody := servicehttp.CalculateRequest{
		Operation: "total_sales",
		StoreID:   "store2",
		StartDate: nil,
		EndDate:   nil,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res servicehttp.CalculateResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res.TotalSales != 0.0 {
		t.Fatalf("expected total sales to be 0.0, got %f", res.TotalSales)
	}
}

func TestIntegration_Calculate_InvalidOperation(t *testing.T) {
	router := setupRouter()

	reqBody := servicehttp.CalculateRequest{
		Operation: "invalid_operation",
		StoreID:   "store1",
		StartDate: nil,
		EndDate:   nil,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res["error"] != "invalid operation" {
		t.Fatalf("expected error message 'invalid operation', got %s", res["error"])
	}
}

func createSale(router *mux.Router, sale *model.Sale, t *testing.T) {
	body, _ := json.Marshal(sale)
	req, err := http.NewRequest("POST", "/data", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
