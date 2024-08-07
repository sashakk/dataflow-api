package sale

import (
	"context"
	"dataflow-api/pkg/model"
	"time"
)

type Service struct {
	saleRepository Repository
}

func NewService(s Repository) *Service {
	return &Service{
		saleRepository: s,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*model.Sale, error) {
	return s.saleRepository.GetAll(ctx)
}

func (s *Service) Create(ctx context.Context, sale *model.Sale) error {
	return s.saleRepository.Create(ctx, sale)
}

func (s *Service) Calculate(ctx context.Context, storeId string, startDate *time.Time, endDate *time.Time) (float64, error) {
	sales, err := s.saleRepository.GetAllByStoreIDAndDateRange(ctx, storeId, startDate, endDate)
	if err != nil {
		return 0, err
	}

	var totalSum float64
	for _, s := range sales {
		totalSum += s.SalePrice
	}

	return totalSum, nil
}
