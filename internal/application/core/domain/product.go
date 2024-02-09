package domain

type Product struct {
	ID             int64
	Name           string `validate:"required"`
	MainCategory   string
	SubCategory    string  `validate:"required"`
	StockNumber    int     `validate:"gte=0"`
	Image          string  `validate:"omitempty,uri"`
	DiscountPrice  float64 `validate:"gte=0"`
	ActualPrice    float64 `validate:"required,gt=0"`
	CurrencyCode   string  `validate:"required,iso4217"`
	CurrencySymbol string
	Version        int64
}

type CreateProductRequest struct {
	Name          string  `validate:"required"`
	SubCategory   string  `validate:"required"`
	StockNumber   int     `validate:"gte=0"`
	Image         string  `validate:"omitempty,uri"`
	DiscountPrice float64 `validate:"gte=0"`
	ActualPrice   float64 `validate:"required,gt=0"`
	CurrencyCode  string  `validate:"required,iso4217"`
}

type UpdateProductRequest struct {
	ID            int64   `validate:"required"`
	Name          string  `validate:"omitempty"`
	SubCategory   string  `validate:"omitempty"`
	StockNumber   int     `validate:"gte=0"`
	Image         string  `validate:"omitempty,uri"`
	DiscountPrice float64 `validate:"gte=0"`
	ActualPrice   float64 `validate:"required,gt=0"`
	CurrencyCode  string  `validate:"omitempty,iso4217"`
	Version       int64   `validate:"gte=1"`
}

type DeleteProductRequest struct {
	ID      int64 `validate:"required"`
	Version int64 `validate:"gte=1"`
}
