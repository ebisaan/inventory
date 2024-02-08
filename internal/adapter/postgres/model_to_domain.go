package postgres

import (
	"github.com/ebisaan/inventory/internal/application/core/domain"
)

func domainProducts(models []*Product) []*domain.Product {
	products := make([]*domain.Product, len(models))
	for i, m := range models {
		products[i] = domainProduct(m)
	}

	return products
}

func domainProduct(model *Product) *domain.Product {
	return &domain.Product{
		ID:             model.ID,
		Name:           model.Name,
		MainCategory:   model.SubCategory.MainCategory.Name,
		SubCategory:    model.SubCategory.Name,
		StockNumber:    model.StockNumber,
		Image:          model.Image,
		DiscountPrice:  model.DiscountPrice,
		ActualPrice:    model.ActualPrice,
		CurrencyCode:   model.Currency.Code,
		CurrencySymbol: model.Currency.Symbol,
		Version:        model.Version,
	}
}
