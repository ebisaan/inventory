package postgres

import "github.com/ebisaan/inventory/internal/application/core/domain"

func insertedProduct(dm *domain.Product) *Product {
	return &Product{
		Name:          dm.Name,
		StockNumber:   dm.StockNumber,
		Image:         dm.Image,
		DiscountPrice: dm.DiscountPrice,
		ActualPrice:   dm.ActualPrice,
		SubCategory: SubCategory{
			Name: dm.SubCategory,
		},
		Currency: Currency{
			Code: dm.CurrencyCode,
		},
	}
}
