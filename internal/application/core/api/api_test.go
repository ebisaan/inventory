package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	api "github.com/ebisaan/inventory/internal/application/core/api"
	"github.com/ebisaan/inventory/internal/application/core/domain"
	"github.com/ebisaan/inventory/internal/application/port"
	mock_port "github.com/ebisaan/inventory/internal/mocks/port"
)

func TestApplication_GetProductByID(t *testing.T) {
	db := mock_port.NewMockDB(t)
	want := &domain.Product{
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
	db.EXPECT().GetProductByID(mock.Anything, int64(1)).Return(want, nil)

	var app port.API
	app, err := api.NewApplication(db)
	require.NoError(t, err)

	got, err := app.GetProductByID(context.Background(), 1)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestApplication_GetProducts(t *testing.T) {
	db := mock_port.NewMockDB(t)
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

	app, err := api.NewApplication(db)
	require.NoError(t, err)

	got, metadata, err := app.GetProducts(context.Background(), domain.Filter{
		Page:     1,
		PageSize: 2,
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

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
	require.NoError(t, err)

	assert.Equal(t, []*domain.Product{want2}, got)
	assert.Equal(t, 1, metadata.PageSize)
	assert.Equal(t, 2, metadata.CurrentPage)
	assert.Equal(t, 1, metadata.FirstPage)
	assert.Equal(t, 2, metadata.LastPage)
	assert.Equal(t, int64(2), metadata.TotalRecords)
}

func TestApplication_CreateProduct(t *testing.T) {
	db := mock_port.NewMockDB(t)
	product := &domain.Product{
		Name:          "Songoku",
		MainCategory:  "Toys & Games",
		SubCategory:   "toys & baby products",
		StockNumber:   10,
		Image:         "",
		DiscountPrice: 0,
		ActualPrice:   50000,
		CurrencyCode:  "VND",
	}
	db.EXPECT().CreateProduct(mock.Anything, product).Return(1, nil)

	var app port.API
	app, err := api.NewApplication(db)
	require.NoError(t, err)

	id, err := app.CreateProduct(context.Background(), product)
	require.NoError(t, err)

	assert.Equal(t, int64(1), id)
}

func TestApplication_CreateProduct_FailedValidation(t *testing.T) {
	db := mock_port.NewMockDB(t)
	product := &domain.Product{
		Name:          "",
		SubCategory:   "",
		StockNumber:   -1,
		Image:         "%notexists$",
		DiscountPrice: -1,
		ActualPrice:   -1,
		CurrencyCode:  "GAY",
	}

	var app port.API
	app, err := api.NewApplication(db)
	require.NoError(t, err)

	_, err = app.CreateProduct(context.Background(), product)
	require.Error(t, err)

	var validationErr domain.ValidationError

	ok := errors.As(err, &validationErr)
	require.True(t, ok)

	assert.Len(t, validationErr.ValidationErrorTranslations, 7)
}
