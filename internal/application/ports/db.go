package port

import (
	"context"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

type DB interface {
	ProductDB
}

type ProductDB interface {
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, filter domain.Filter) (int64, []*domain.Product, error)
}
