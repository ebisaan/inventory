package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	"github.com/ebisaan/inventory/internal/application/port"
)

var _ port.DB = (*Adapter)(nil)

type Adapter struct {
	db *gorm.DB
}

type Config struct {
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
}

func NewAdapter(dsn string, cfg ...Config) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		FullSaveAssociations: false,
		// Logger:                 nil,
		DisableAutomaticPing: false,
		TranslateError:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("open gorm connection pool: %w", err)
	}

	if len(cfg) > 0 {
		cfg := cfg[0]
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("get *sql.DB: %w", err)
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)
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

func (a *Adapter) CreateProduct(ctx context.Context, dp *domain.CreateProductRequest) (id int64, err error) {
	db := a.db.WithContext(ctx)

	p := insertedProduct(dp)

	tx := db.Begin()
	defer func() {
		var txErr error
		if err == nil {
			txErr = tx.Commit().Error
		} else {
			txErr = tx.Rollback().Error
		}

		if txErr != nil {
			err = fmt.Errorf("%w: %w", txErr, err)
		}
	}()

	scID, err := getSubcategoryIDByName(tx, p.SubCategory.Name)
	if err != nil {
		return 0, fmt.Errorf("select subcategory id: %w", err)
	}
	p.SubCategoryID = scID

	crcID, err := getCurrencyIDByCode(tx, p.Currency.Code)
	if err != nil {
		return 0, fmt.Errorf("select currency id: %w", err)
	}
	p.CurrencyID = crcID

	err = tx.Omit(clause.Associations).Create(&p).Error
	if err != nil {
		return 0, fmt.Errorf("insert product: %w", err)
	}

	return p.ID, nil
}

func (a *Adapter) UpdateProduct(ctx context.Context, id int64, dp *domain.UpdateProductRequest) (err error) {
	db := a.db.WithContext(ctx)

	p := updatedProduct(dp)
	p.ID = id
	curVersion := p.Version
	p.Version += 1

	tx := db.Begin()
	defer func() {
		var txErr error
		if err == nil {
			txErr = tx.Commit().Error
		} else {
			txErr = tx.Rollback().Error
		}

		if txErr != nil {
			err = fmt.Errorf("%w: %w", txErr, err)
		}
	}()

	scID, err := getSubcategoryIDByName(tx, p.SubCategory.Name)
	if err != nil {
		return fmt.Errorf("select subcategory id: %w", err)
	}
	p.SubCategoryID = scID

	crcID, err := getCurrencyIDByCode(tx, p.Currency.Code)
	if err != nil {
		return fmt.Errorf("select currency id: %w", err)
	}
	p.CurrencyID = crcID

	res := tx.Omit(clause.Associations).Where("id = ?", id).Where("version = ?", curVersion).Updates(&p)
	if err := res.Error; err != nil {
		return fmt.Errorf("select product by id=%d: %w", id, err)
	}

	if res.RowsAffected == 0 {
		return domain.ErrEditConflict
	}

	return nil
}

func (a *Adapter) IsSubCategoryExists(ctx context.Context, name string) (bool, error) {
	var found bool
	err := a.db.
		Model(&SubCategory{}).
		Select("count(*) > 0").
		Where("name = ?", name).
		Take(&found).
		Error
	if err != nil {
		return false, err
	}

	return found, nil
}

func (a *Adapter) IsCurrencyCodeExists(ctx context.Context, code string) (bool, error) {
	var found bool
	err := a.db.
		Model(&Currency{}).
		Select("count(*) > 0").
		Where("code = ?", code).
		Take(&found).
		Error
	if err != nil {
		return false, err
	}

	return found, nil
}

func getSubcategoryIDByName(db *gorm.DB, name string) (int64, error) {
	var id int64
	err := db.Model(&SubCategory{}).Select("id").Where("name = ?", name).First(&id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return 0, domain.ErrAssociationNotFound
		default:
			return 0, err
		}
	}

	return id, nil
}

func getCurrencyIDByCode(db *gorm.DB, code string) (int64, error) {
	var id int64
	err := db.Model(&Currency{}).Select("id").Where("code = ?", code).First(&id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return 0, domain.ErrAssociationNotFound
		default:
			return 0, fmt.Errorf("select currency id: %w", err)
		}
	}

	return id, nil
}

func (a *Adapter) AutoMigration(ctx context.Context) error {
	db := a.db.WithContext(ctx)
	err := db.AutoMigrate(&Currency{}, &MainCategory{}, &SubCategory{}, &Product{})
	if err != nil {
		return fmt.Errorf("auto migration: %w", err)
	}

	return nil
}
