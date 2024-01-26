package postgres

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dsn string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		FullSaveAssociations: false,
		// Logger:                 nil,
		DisableAutomaticPing: false,
		TranslateError:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("open gorm connection pool: %w", err)
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) InsertProduct(ctx context.Context, domainProduct *domain.Product) (*domain.Product, error) {
	return nil, nil
}

func (a *Adapter) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	db := a.db.WithContext(ctx)
	product := &Product{}
	err := db.Joins("SubCategory.MainCategory").Joins("Currency").First(product, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.ErrNotFound
		default:
			return nil, fmt.Errorf("select product by id=%d: %w", id, err)
		}
	}

	return domainProduct(product), nil
}

func (a *Adapter) GetProducts(ctx context.Context, filter domain.Filter) (int64, []*domain.Product, error) {
	db := a.db.WithContext(ctx)

	var products []*Product
	query := db.Model(&products)

	var total int64 = 0
	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, fmt.Errorf("count products: %w", err)
	}
	if total > 0 {
		err := query.Joins("SubCategory.MainCategory").Joins("Currency").
			Limit(filter.Limit()).
			Offset(int(filter.Offset())).
			Find(&products).
			Error
		if err != nil {
			return 0, nil, fmt.Errorf("select products: %w", err)
		}
	}

	return total, domainProducts(products), nil
}

func (a *Adapter) AutoMigration(ctx context.Context) error {
	db := a.db.WithContext(ctx)
	err := db.AutoMigrate(&Currency{}, &MainCategory{}, &SubCategory{}, &Product{})
	if err != nil {
		return fmt.Errorf("auto migration: %w", err)
	}

	return nil
}
