package api

import (
	"context"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	port "github.com/ebisaan/inventory/internal/application/ports"
)

type Product struct {
	db port.DB
}

func (p *Product) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return p.db.GetProductByID(ctx, id)
}

func (p *Product) GetAllProducts(ctx context.Context, filter domain.Filter) ([]*domain.Product, domain.Metadata, error) {
	return nil, domain.Metadata{}, nil
}
