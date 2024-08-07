package sale

import (
	"context"
	"dataflow-api/pkg/model"
	"time"
)

type Repository interface {
	Create(context.Context, *model.Sale) error
	GetAllByStoreIDAndDateRange(context.Context, string, *time.Time, *time.Time) ([]*model.Sale, error)
	GetAll(context.Context) ([]*model.Sale, error)
}
