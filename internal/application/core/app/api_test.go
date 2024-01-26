package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ebisaan/inventory/internal/application/core/app"
	"github.com/ebisaan/inventory/internal/application/core/domain"
	"github.com/ebisaan/inventory/internal/mocks/port"
)

func TestApplicationGetProductByID(t *testing.T) {
	db := port.NewMockDB(t)
	want := &domain.Product{
		ID:             0,
		Name:           "Songoku",
		MainCategory:   "Toys & Games",
		SubCategory:    "toys & baby products",
		StockNumber:    10,
		Image:          "",
		DiscountPrice:  0,
		ActualPrice:    50000,
		CurrencyCode:   "VND",
		CurrencySymbol: "₫",
	}
	db.EXPECT().GetProductByID(mock.Anything, int64(1)).Return(want, nil)

	app := app.NewApplication(db)

	got, err := app.GetProductByID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestApplicationGetProducts(t *testing.T) {
	db := port.NewMockDB(t)
	want1 := &domain.Product{
		ID:             1,
		Name:           "Songoku",
		MainCategory:   "Toys & Games",
		SubCategory:    "toys & baby products",
		StockNumber:    10,
		Image:          "",
		DiscountPrice:  0,
		ActualPrice:    50000,
		CurrencyCode:   "VND",
		CurrencySymbol: "₫",
	}

	want2 := &domain.Product{
		ID:             2,
		Name:           "G-Shock",
		MainCategory:   "Toys & Games",
		SubCategory:    "toys & baby products",
		StockNumber:    100,
		Image:          "",
		DiscountPrice:  0,
		ActualPrice:    500,
		CurrencyCode:   "USD",
		CurrencySymbol: "$",
	}

	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     1,
		PageSize: 2,
	}).Return(int64(2), []*domain.Product{want1, want2}, nil)
	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     1,
		PageSize: 1,
	}).Return(int64(2), []*domain.Product{want1}, nil)
	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     2,
		PageSize: 1,
	}).Return(int64(2), []*domain.Product{want2}, nil)

	app := app.NewApplication(db)

	got, metadata, err := app.GetProducts(context.Background(), domain.Filter{
		Page:     1,
		PageSize: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*domain.Product{want1, want2}, got)
	assert.Equal(t, 2, metadata.PageSize)
	assert.Equal(t, 1, metadata.CurrentPage)
	assert.Equal(t, 1, metadata.FirstPage)
	assert.Equal(t, 1, metadata.LastPage)
	assert.Equal(t, int64(2), metadata.TotalRecords)

	got, metadata, err = app.GetProducts(context.Background(), domain.Filter{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*domain.Product{want1}, got)
	assert.Equal(t, 1, metadata.PageSize)
	assert.Equal(t, 1, metadata.CurrentPage)
	assert.Equal(t, 1, metadata.FirstPage)
	assert.Equal(t, 2, metadata.LastPage)
	assert.Equal(t, int64(2), metadata.TotalRecords)

	got, metadata, err = app.GetProducts(context.Background(), domain.Filter{
		Page:     2,
		PageSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*domain.Product{want2}, got)
	assert.Equal(t, 1, metadata.PageSize)
	assert.Equal(t, 2, metadata.CurrentPage)
	assert.Equal(t, 1, metadata.FirstPage)
	assert.Equal(t, 2, metadata.LastPage)
	assert.Equal(t, int64(2), metadata.TotalRecords)
}
