package port

import (
	"context"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

type API interface {
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, filter domain.Filter) ([]*domain.Product, domain.Metadata, error)
	CreateProduct(ctx context.Context, product *domain.Product) (id int64, err error)
}
