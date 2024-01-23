package domain

type Product struct {
	ID             int64
	Name           string
	MainCategory   string
	SubCategory    string
	StockNumber    int
	Image          string
	DiscountPrice  float64
	ActualPrice    float64
	CurrencyCode   string
	CurrencySymbol string
}
