package sale_test

import (
	"context"
	"dataflow-api/pkg/model"
	"dataflow-api/pkg/repository/sale"
	"testing"
	"time"
)

func TestInMemorySaleRepository_Create(t *testing.T) {
	repo := sale.NewInMemorySaleRepository()
	sale := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     time.Now(),
	}

	err := repo.Create(context.Background(), sale)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	sales, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sales) != 1 {
		t.Fatalf("expected 1 sale, got %d", len(sales))
	}

	if sales[0].ProductID != "p1" {
		t.Fatalf("expected ProductID to be 'p1', got %s", sales[0].ProductID)
	}
}

func TestInMemorySaleRepository_GetAll(t *testing.T) {
	repo := sale.NewInMemorySaleRepository()
	sale1 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     time.Now(),
	}
	sale2 := &model.Sale{
		ProductID:    "p2",
		StoreID:      "store2",
		QuantitySold: 5,
		SalePrice:    200.0,
		SaleDate:     time.Now(),
	}
	repo.Create(context.Background(), sale1)
	repo.Create(context.Background(), sale2)

	sales, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sales) != 2 {
		t.Fatalf("expected 2 sales, got %d", len(sales))
	}
}
