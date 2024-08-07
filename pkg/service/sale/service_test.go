package sale_test

import (
	"context"
	"dataflow-api/pkg/model"
	repositorysale "dataflow-api/pkg/repository/sale"
	"dataflow-api/pkg/service/sale"
	"testing"
	"time"
)

func createWithoutErr(t *testing.T, service *sale.Service, s *model.Sale) {
	err := service.Create(context.Background(), s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestService_Create(t *testing.T) {
	repo := repositorysale.NewInMemorySaleRepository()
	service := sale.NewService(repo)

	s := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     time.Now(),
	}

	createWithoutErr(t, service, s)

	sales, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sales) != 1 {
		t.Fatalf("expected 1 sale, got %d", len(sales))
	}

	if sales[0].ProductID != s.ProductID {
		t.Fatalf("expected ProductID to be '%s', got '%s'", s.ProductID, sales[0].ProductID)
	}
}

func TestService_GetAll(t *testing.T) {
	repo := repositorysale.NewInMemorySaleRepository()
	service := sale.NewService(repo)

	s1 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     time.Now(),
	}
	s2 := &model.Sale{
		ProductID:    "p2",
		StoreID:      "store2",
		QuantitySold: 5,
		SalePrice:    200.0,
		SaleDate:     time.Now(),
	}

	createWithoutErr(t, service, s1)
	createWithoutErr(t, service, s2)

	sales, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sales) != 2 {
		t.Fatalf("expected 2 sales, got %d", len(sales))
	}
}

func TestService_Calculate(t *testing.T) {
	repo := repositorysale.NewInMemorySaleRepository()
	service := sale.NewService(repo)

	now := time.Now()
	s1 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     now.AddDate(0, 0, -1),
	}
	s2 := &model.Sale{
		ProductID:    "p2",
		StoreID:      "store1",
		QuantitySold: 5,
		SalePrice:    200.0,
		SaleDate:     now,
	}
	s3 := &model.Sale{
		ProductID:    "p3",
		StoreID:      "store2",
		QuantitySold: 3,
		SalePrice:    300.0,
		SaleDate:     now,
	}

	createWithoutErr(t, service, s1)
	createWithoutErr(t, service, s2)
	createWithoutErr(t, service, s3)

	startDate := now.AddDate(0, 0, -7)
	endDate := now

	total, err := service.Calculate(context.Background(), "store1", &startDate, &endDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTotal := 300.0
	if total != expectedTotal {
		t.Fatalf("expected total %f, got %f", expectedTotal, total)
	}
}

func TestService_Calculate_NoDateRange(t *testing.T) {
	repo := repositorysale.NewInMemorySaleRepository()
	service := sale.NewService(repo)

	now := time.Now()
	s1 := &model.Sale{
		ProductID:    "p1",
		StoreID:      "store1",
		QuantitySold: 10,
		SalePrice:    100.0,
		SaleDate:     now.AddDate(0, 0, -10),
	}
	s2 := &model.Sale{
		ProductID:    "p2",
		StoreID:      "store1",
		QuantitySold: 5,
		SalePrice:    200.0,
		SaleDate:     now,
	}

	createWithoutErr(t, service, s1)
	createWithoutErr(t, service, s2)

	total, err := service.Calculate(context.Background(), "store1", nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTotal := 300.0
	if total != expectedTotal {
		t.Fatalf("expected total %f, got %f", expectedTotal, total)
	}
}
