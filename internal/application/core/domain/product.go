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
}
