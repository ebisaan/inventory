package postgres

import (
	"time"
)

type BaseModel struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Product struct {
	BaseModel

	Name        string `gorm:"notnull"`
	StockNumber int    `gorm:"type=integer;notnull;check:stock_number > 0"`
	Image       string

	DiscountPrice float64 `gorm:"check:discount_price >= 0"`
	ActualPrice   float64 `gorm:"notnull;check:actual_price > 0"`

	SubCategoryID int64 `gorm:"notnull"`
	SubCategory   SubCategory

	CurrencyID int64 `gorm:"notnull"`
	Currency   Currency
}

type MainCategory struct {
	BaseModel
	Name string `gorm:"notnull;uniqueIndex"`
}

type SubCategory struct {
	BaseModel
	Name           string `gorm:"notnull;uniqueIndex"`
	MainCategoryID int64  `gorm:"notnull"`
	MainCategory   MainCategory
}

type Currency struct {
	BaseModel
	Code   string `gorm:"notnull;uniqueIndex"`
	Symbol string `gorm:"notnull;uniqueIndex"`
}
