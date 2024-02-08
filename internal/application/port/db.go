package port

import (
	"context"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

type DB interface {
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, filter domain.Filter) (int64, []*domain.Product, error)
	CreateProduct(ctx context.Context, dp *domain.CreateProductRequest) (id int64, err error)
	UpdateProduct(ctx context.Context, id int64, db *domain.UpdateProductRequest) error
	IsSubCategoryExists(ctx context.Context, subCategory string) (bool, error)
	IsCurrencyCodeExists(ctx context.Context, currencyCode string) (bool, error)
}
