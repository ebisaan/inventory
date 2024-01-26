package app

import (
	"context"
	"fmt"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	port "github.com/ebisaan/inventory/internal/application/port"
)

type Application struct {
	db port.DB
}

func NewApplication(db port.DB) *Application {
	return &Application{
		db: db,
	}
}

func (a *Application) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return a.db.GetProductByID(ctx, id)
}

func (a *Application) GetProducts(ctx context.Context, filter domain.Filter) ([]*domain.Product, domain.Metadata, error) {
	if filter.PageSize == 0 {
		filter.PageSize = domain.DefaultPageSize
	}
	if filter.PageSize > domain.MaxPageSize {
		filter.PageSize = domain.MaxPageSize
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	n, products, err := a.db.GetProducts(ctx, filter)
	if err != nil {
		return nil, domain.Metadata{}, fmt.Errorf("get products from db: %w", err)
	}

	metadata := domain.MakeMetadata(n, filter.Page, filter.PageSize)

	return products, metadata, nil
}
