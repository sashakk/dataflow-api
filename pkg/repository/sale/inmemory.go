package sale

import (
	"context"
	"dataflow-api/pkg/model"
	"sync"
	"time"
)

type InMemorySaleRepository struct {
	db map[string][]*model.Sale
	mu sync.RWMutex
}

func NewInMemorySaleRepository() *InMemorySaleRepository {
	return &InMemorySaleRepository{db: make(map[string][]*model.Sale)}
}

func (r *InMemorySaleRepository) Create(_ context.Context, sale *model.Sale) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[sale.StoreID] = append(r.db[sale.StoreID], sale)
	return nil
}

func (r *InMemorySaleRepository) GetAll(_ context.Context) ([]*model.Sale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allSales := make([]*model.Sale, 0)
	for _, sales := range r.db {
		allSales = append(allSales, sales...)
	}

	return allSales, nil
}

func (r *InMemorySaleRepository) GetAllByStoreIDAndDateRange(_ context.Context, storeId string, startDate *time.Time, endDate *time.Time) ([]*model.Sale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Sale
	sales, ok := r.db[storeId]
	if !ok {
		return result, nil
	}

	for _, s := range sales {
		if (startDate == nil || !s.SaleDate.Before(*startDate)) &&
			(endDate == nil || !s.SaleDate.After(*endDate)) {
			result = append(result, s)
		}
	}

	return result, nil
}
