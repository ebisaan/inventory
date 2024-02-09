package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	progresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	port "github.com/ebisaan/inventory/internal/application/port"
)

type DatabaseTestSuite struct {
	suite.Suite
	db             port.DB
	container      testcontainers.Container
	dsn            string
	domainProducts []*domain.Product
	products       []*Product
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) TestGetProductByID() {
	p, err := s.db.GetProductByID(context.Background(), s.domainProducts[0].ID)
	s.Require().NoError(err)

	s.Suite.Assert().Equal(s.domainProducts[0], p)
}

func (s *DatabaseTestSuite) TestGetProductByID_NotFound() {
	_, err := s.db.GetProductByID(context.Background(), s.domainProducts[len(s.domainProducts)-1].ID+1)
	s.Require().Error(err)
	if !errors.Is(err, domain.ErrNotFound) {
		s.T().Errorf("got error %q, want %q", err, domain.ErrNotFound)
	}
}

func (s *DatabaseTestSuite) TestGetProducts() {
	s.Run("get all products", func() {
		n, products, err := s.db.GetProducts(context.Background(), domain.Filter{
			Page:     1,
			PageSize: domain.DefaultPageSize,
		})
		if err != nil {
			s.T().Fatal(err)
		}

		s.Suite.Assert().Equal(s.domainProducts, products)
		s.Suite.Assert().Equal(int64(2), n)
	})

	s.Run("get products by page", func() {
		n, products, err := s.db.GetProducts(context.Background(), domain.Filter{
			Page:     1,
			PageSize: 1,
		})
		if err != nil {
			s.T().Fatal(err)
		}
		s.Suite.Assert().Equal([]*domain.Product{s.domainProducts[0]}, products)
		s.Suite.Assert().Equal(int64(2), n)

		n, products, err = s.db.GetProducts(context.Background(), domain.Filter{
			Page:     2,
			PageSize: 1,
		})
		if err != nil {
			s.T().Fatal(err)
		}

		s.Suite.Assert().Equal([]*domain.Product{s.domainProducts[1]}, products)
		s.Suite.Assert().Equal(int64(2), n)

		n, products, err = s.db.GetProducts(context.Background(), domain.Filter{
			Page:     3,
			PageSize: 1,
		})
		if err != nil {
			s.T().Fatal(err)
		}

		s.Suite.Assert().Equal([]*domain.Product{}, products)
		s.Suite.Assert().Equal(int64(2), n)
	})
}

func (s *DatabaseTestSuite) TestCreateProduct() {
	createRequest := &domain.CreateProductRequest{
		Name:          "Superman",
		SubCategory:   "Toys & Games",
		StockNumber:   100,
		Image:         "image.com/123",
		DiscountPrice: 0,
		ActualPrice:   10,
		CurrencyCode:  "VND",
	}

	ctx := context.Background()
	id, err := s.db.CreateProduct(ctx, createRequest)
	s.Require().NoError(err)

	db := s.getGormDB()

	var gotProduct Product
	err = db.Preload(clause.Associations).First(&gotProduct, id).Error
	s.Require().NoError(err)
	s.Assert().NotNil(gotProduct)
	s.Assert().Equal(createRequest.Name, gotProduct.Name)
	s.Assert().Equal(createRequest.SubCategory, gotProduct.SubCategory.Name)
	s.Assert().Equal(createRequest.CurrencyCode, gotProduct.Currency.Code)
	s.Assert().Equal(createRequest.StockNumber, gotProduct.StockNumber)
	s.Assert().Equal(createRequest.Image, gotProduct.Image)
	s.Assert().Equal(createRequest.DiscountPrice, gotProduct.DiscountPrice)
	s.Assert().Equal(createRequest.ActualPrice, gotProduct.ActualPrice)
	s.Assert().Equal(int64(1), gotProduct.Version)

	err = db.Delete(&gotProduct).Error
	s.Require().NoError(err)
}

func (s *DatabaseTestSuite) TestUpdateProduct() {
	p := Product{
		Name:          "Superman",
		SubCategory:   s.products[0].SubCategory,
		Currency:      s.products[0].Currency,
		StockNumber:   100,
		Image:         "image.com/123",
		DiscountPrice: 0,
		ActualPrice:   10,
		Version:       1,
	}

	db := s.getGormDB()

	err := db.Save(&p).Error
	s.Require().NoError(err)
	s.Assert().NotEmpty(p.ID)

	updateRequest := &domain.UpdateProductRequest{
		ID:            p.ID,
		Name:          "Fake Superman",
		SubCategory:   s.products[1].SubCategory.Name,
		CurrencyCode:  s.products[1].Currency.Code,
		StockNumber:   10,
		Image:         "image.com/456",
		DiscountPrice: 10,
		ActualPrice:   100,
		Version:       1,
	}

	err = s.db.UpdateProduct(context.Background(), updateRequest)
	s.Require().NoError(err)

	gotProduct := &Product{}
	err = db.Preload(clause.Associations).First(&gotProduct, p.ID).Error
	s.Require().NoError(err)

	s.Assert().Equal(updateRequest.Name, gotProduct.Name)
	s.Assert().Equal(updateRequest.SubCategory, gotProduct.SubCategory.Name)
	s.Assert().Equal(updateRequest.CurrencyCode, gotProduct.Currency.Code)
	s.Assert().Equal(updateRequest.StockNumber, gotProduct.StockNumber)
	s.Assert().Equal(updateRequest.Image, gotProduct.Image)
	s.Assert().Equal(updateRequest.DiscountPrice, gotProduct.DiscountPrice)
	s.Assert().Equal(updateRequest.ActualPrice, gotProduct.ActualPrice)
	s.Assert().Equal(updateRequest.Version+1, gotProduct.Version)
	err = db.Delete(&Product{}, p.ID).Error
	s.Require().NoError(err)
}

func (s *DatabaseTestSuite) TestUpdateProduct_Conflict() {
	p := Product{
		Name:          "Superman",
		SubCategory:   s.products[0].SubCategory,
		Currency:      s.products[0].Currency,
		StockNumber:   100,
		Image:         "image.com/123",
		DiscountPrice: 0,
		ActualPrice:   10,
		Version:       2,
	}

	db := s.getGormDB()

	err := db.Save(&p).Error
	s.Require().NoError(err)
	s.Assert().NotEmpty(p.ID)

	updateRequest := &domain.UpdateProductRequest{
		ID:            p.ID,
		Name:          "Fake Superman",
		SubCategory:   s.products[1].SubCategory.Name,
		CurrencyCode:  s.products[1].Currency.Code,
		StockNumber:   10,
		Image:         "image.com/456",
		DiscountPrice: 10,
		ActualPrice:   100,
		Version:       1,
	}

	err = s.db.UpdateProduct(context.Background(), updateRequest)
	s.Require().Error(err)
	if !errors.Is(err, domain.ErrEditConflict) {
		s.T().Errorf("got error %q, want %q", err, domain.ErrEditConflict)
	}

	err = db.Delete(&Product{}, p.ID).Error
	s.Require().NoError(err)
}

func (s *DatabaseTestSuite) TestDeleteProduct() {
	p := Product{
		Name:          "Superman",
		SubCategory:   s.products[0].SubCategory,
		Currency:      s.products[0].Currency,
		StockNumber:   100,
		Image:         "image.com/123",
		DiscountPrice: 0,
		ActualPrice:   10,
		Version:       1,
	}

	db := s.getGormDB()

	err := db.Save(&p).Error
	s.Require().NoError(err)
	s.Assert().NotEmpty(p.ID)

	req := &domain.DeleteProductRequest{
		ID:      p.ID,
		Version: p.Version,
	}

	err = s.db.DeleteProduct(context.Background(), req)
	s.Require().NoError(err)
}

func (s *DatabaseTestSuite) SetupSuite() {
	s.setupContainer()
	s.setupAdapter()
	s.setupProducts()
}

func (s *DatabaseTestSuite) getGormDB() *gorm.DB {
	s.T().Helper()
	db, err := gorm.Open(progresDriver.Open(s.dsn), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("open gorm connection pool: %s", err)
	}

	return db
}

func (s *DatabaseTestSuite) setupProducts() {
	db, err := gorm.Open(progresDriver.Open(s.dsn), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("open gorm connection pool: %s", err)
	}

	products := []*Product{
		{
			Name:        "Songoku",
			StockNumber: 10,
			ActualPrice: 500000,
			SubCategory: SubCategory{
				Name: "Toys & Games",
				MainCategory: MainCategory{
					Name: "toys & baby products",
				},
			},
			Currency: Currency{
				Code:   "VND",
				Symbol: "₫",
			},
		},
		{
			Name:        "G-Shock",
			StockNumber: 100,
			ActualPrice: 500,
			SubCategory: SubCategory{
				Name: "Watches",
				MainCategory: MainCategory{
					Name: "accessories",
				},
			},
			Currency: Currency{
				Code:   "USD",
				Symbol: "$",
			},
		},
	}

	err = db.Create(&products).Error
	if err != nil {
		s.T().Fatalf("create products: %s", err)
	}

	s.domainProducts = domainProducts(products)
	s.products = products
}

func (s *DatabaseTestSuite) setupAdapter() {
	ctx := context.Background()

	postgresDB, err := NewAdapter(s.dsn)
	if err != nil {
		s.T().Fatalf("create new postgres adapter: %s", err)
	}

	err = postgresDB.AutoMigration(ctx)
	if err != nil {
		s.T().Fatalf("auto migration: %s", err)
	}
	s.db = postgresDB
}

func (s *DatabaseTestSuite) setupContainer() {
	ctx := context.Background()

	user := "test-user"
	password := "secret-password"
	exposePort := "5432/tcp"
	dbName := "ebisaan"
	url := func(address string) string {
		return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, address, dbName)
	}
	dbURL := func(host string, port nat.Port) string {
		return url(fmt.Sprintf("%s:%s", host, port.Port()))
	}

	req := testcontainers.ContainerRequest{
		Image:        "docker.io/postgres:16-alpine",
		ExposedPorts: []string{exposePort},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForSQL(nat.Port(exposePort), "postgres", dbURL),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		s.T().Fatalf("set up postgreSQL container: %s", err)
	}
	s.container = postgresContainer
	ctnHost, _ := s.container.Endpoint(ctx, "")
	s.dsn = url(ctnHost)
}

func (s *DatabaseTestSuite) TeardownSuite() {
	err := s.container.Terminate(context.Background())
	if err != nil {
		s.T().Fatalf("terminate postgreSQL container: %s", err)
	}
}
