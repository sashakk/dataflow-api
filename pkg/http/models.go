package http

import "time"

const (
	operationTotalSales = "total_sales"
)

type CalculateRequest struct {
	Operation string     `json:"operation"`
	StoreID   string     `json:"store_id"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type CalculateResponse struct {
	StoreID    string     `json:"store_id"`
	TotalSales float64    `json:"total_sales"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
}

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type successResponse struct {
	Status string `json:"status"`
}
