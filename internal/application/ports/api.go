package port

import (
	"context"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

type API interface {
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetAllProducts(ctx context.Context) ([]*domain.Product, error)
}
