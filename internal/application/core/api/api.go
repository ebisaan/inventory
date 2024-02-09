package api

import (
	"context"
	"fmt"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	port "github.com/ebisaan/inventory/internal/application/port"
)

var _ port.API = (*Application)(nil)

type Application struct {
	db port.DB
	v  *validate
}

func NewApplication(db port.DB) (*Application, error) {
	v, err := newValidate("json")
	if err != nil {
		return nil, err
	}
	return &Application{
		db: db,
		v:  v,
	}, nil
}

func (a *Application) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return a.db.GetProductByID(ctx, id)
}

func (a *Application) GetProducts(ctx context.Context, filter domain.Filter) ([]*domain.Product, domain.Metadata, error) {
	filter = domain.ProcessFilter(filter)
	n, products, err := a.db.GetProducts(ctx, filter)
	if err != nil {
		return nil, domain.Metadata{}, fmt.Errorf("get products from db: %w", err)
	}

	metadata := domain.MakeMetadata(n, filter.Page, filter.PageSize)

	return products, metadata, nil
}

func (a *Application) CreateProduct(ctx context.Context, product *domain.CreateProductRequest) (id int64, err error) {
	err = a.v.ValidateStruct(product)
	if err != nil {
		return 0, err
	}

	found, err := a.db.IsSubCategoryExists(ctx, product.SubCategory)
	if err != nil {
		return 0, fmt.Errorf("is subcategory exists: %w", err)
	}
	if !found {
		return 0, domain.ValidationError{
			FieldErrorMessages: map[string]string{
				"Subcategory": "not exists",
			},
		}
	}

	found, err = a.db.IsCurrencyCodeExists(ctx, product.CurrencyCode)
	if err != nil {
		return 0, fmt.Errorf("is currency code exists: %w", err)
	}

	if !found {
		return 0, domain.ValidationError{
			FieldErrorMessages: map[string]string{
				"Subcategory": "not exists",
			},
		}
	}

	id, err = a.db.CreateProduct(ctx, product)
	if err != nil {
		return 0, fmt.Errorf("create product: %w", err)
	}

	return id, nil
}

func (a *Application) UpdateProduct(ctx context.Context, req *domain.UpdateProductRequest) error {
	err := a.v.ValidateStruct(req)
	if err != nil {
		return err
	}

	err = a.db.UpdateProduct(ctx, req)
	return err
}

func (a *Application) DeleteProduct(ctx context.Context, req *domain.DeleteProductRequest) error {
	err := a.v.ValidateStruct(req)
	if err != nil {
		return err
	}

	err = a.db.DeleteProduct(ctx, req)
	return err
}
