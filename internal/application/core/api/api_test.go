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

var readOnlyTestProducts = [...]*domain.Product{
	{
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
	},
	{
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
	},
}

func TestApplication_GetProductByID(t *testing.T) {
	db := mock_port.NewMockDB(t)
	want := readOnlyTestProducts[0]
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

	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     1,
		PageSize: 2,
	}).Return(int64(2), readOnlyTestProducts[:], nil)
	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     1,
		PageSize: 1,
	}).Return(int64(2), readOnlyTestProducts[:1], nil)
	db.EXPECT().GetProducts(mock.Anything, domain.Filter{
		Page:     2,
		PageSize: 1,
	}).Return(int64(2), readOnlyTestProducts[1:2], nil)

	app, err := api.NewApplication(db)
	require.NoError(t, err)

	got, metadata, err := app.GetProducts(context.Background(), domain.Filter{
		Page:     1,
		PageSize: 2,
	})
	require.NoError(t, err)

	assert.Equal(t, readOnlyTestProducts[:], got)
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

	assert.Equal(t, readOnlyTestProducts[:1], got)
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

	assert.Equal(t, readOnlyTestProducts[1:2], got)
	assert.Equal(t, 1, metadata.PageSize)
	assert.Equal(t, 2, metadata.CurrentPage)
	assert.Equal(t, 1, metadata.FirstPage)
	assert.Equal(t, 2, metadata.LastPage)
	assert.Equal(t, int64(2), metadata.TotalRecords)
}

func TestApplication_CreateProduct(t *testing.T) {
	db := mock_port.NewMockDB(t)
	product := &domain.CreateProductRequest{
		Name:          "Songoku",
		SubCategory:   "toys & baby products",
		StockNumber:   10,
		Image:         "",
		DiscountPrice: 0,
		ActualPrice:   50000,
		CurrencyCode:  "VND",
	}
	db.EXPECT().IsCurrencyCodeExists(mock.Anything, "VND").Return(true, nil)
	db.EXPECT().IsSubCategoryExists(mock.Anything, "toys & baby products").Return(true, nil)
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
	product := &domain.CreateProductRequest{
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

	assert.Len(t, validationErr.FieldErrorMessages, 7)
}

func TestApplication_UpdateProduct(t *testing.T) {
	db := mock_port.NewMockDB(t)
	product := &domain.UpdateProductRequest{
		Name:          "Songoku",
		SubCategory:   "toys & baby products",
		StockNumber:   10,
		Image:         "",
		DiscountPrice: 0,
		ActualPrice:   50000,
		CurrencyCode:  "VND",
		Version:       1,
	}
	var id int64 = 1
	db.EXPECT().UpdateProduct(mock.Anything, id, product).Return(nil)

	var app port.API
	app, err := api.NewApplication(db)
	require.NoError(t, err)

	err = app.UpdateProduct(context.Background(), id, product)
	require.NoError(t, err)
}

func TestApplication_UpdateProduct_FailedValidation(t *testing.T) {
	db := mock_port.NewMockDB(t)
	product := &domain.UpdateProductRequest{
		StockNumber:   -1,
		Image:         "%notexists$",
		DiscountPrice: -1,
		ActualPrice:   -1,
		CurrencyCode:  "GAY",
		Version:       -1,
	}

	var app port.API
	app, err := api.NewApplication(db)
	require.NoError(t, err)

	err = app.UpdateProduct(context.Background(), 1, product)
	require.Error(t, err)

	var validationErr domain.ValidationError

	ok := errors.As(err, &validationErr)
	require.True(t, ok)

	assert.Len(t, validationErr.FieldErrorMessages, 6)
}
