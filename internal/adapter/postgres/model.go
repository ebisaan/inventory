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

	Name        string `gorm:"not null"`
	StockNumber int    `gorm:"type=integer;not null;check:stock_number > 0"`
	Image       string

	DiscountPrice float64 `gorm:"check:discount_price >= 0"`
	ActualPrice   float64 `gorm:"not null;check:actual_price >= 0"`

	SubCategoryID int64 `gorm:"not null"`
	SubCategory   SubCategory

	CurrencyID int64 `gorm:"not null"`
	Currency   Currency

	Version int64 `gorm:"not null;default:1"`
}

type MainCategory struct {
	BaseModel
	Name string `gorm:"not null;uniqueIndex"`
}

type SubCategory struct {
	BaseModel
	Name           string `gorm:"not null;uniqueIndex"`
	MainCategoryID int64  `gorm:"not null"`
	MainCategory   MainCategory
}

type Currency struct {
	BaseModel
	Code   string `gorm:"not null;uniqueIndex"`
	Symbol string `gorm:"not null;uniqueIndex"`
}
