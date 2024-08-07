package model

import (
	"fmt"
	"time"
)

type Sale struct {
	ProductID    string    `json:"product_id"`
	StoreID      string    `json:"store_id"`
	QuantitySold int       `json:"quantity_sold"`
	SalePrice    float64   `json:"sale_price"`
	SaleDate     time.Time `json:"sale_date"`
}

func ValidateSale(s *Sale) error {
	if s.ProductID == "" {
		return fmt.Errorf("product_id is required")
	}
	if s.StoreID == "" {
		return fmt.Errorf("store_id is required")
	}
	if s.QuantitySold <= 0 {
		return fmt.Errorf("quantity_sold must be greater than 0")
	}
	if s.SalePrice <= 0 {
		return fmt.Errorf("sale_price must be greater than 0")
	}
	if s.SaleDate.IsZero() {
		return fmt.Errorf("sale_date is required")
	}

	return nil
}
