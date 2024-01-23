package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	progresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ebisaan/inventory/internal/application/core/domain"
	port "github.com/ebisaan/inventory/internal/application/ports"
)

type DatabaseTestSuite struct {
	suite.Suite
	db        port.DB
	container testcontainers.Container
	dsn       string
	products  []*domain.Product
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) TestGetProductByID() {
	p, err := s.db.GetProductByID(context.Background(), s.products[0].ID)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Suite.Assert().Equal(s.products[0], p)
}

func (s *DatabaseTestSuite) TestGetProducts() {
	s.Run("get all products", func() {
		n, products, err := s.db.GetProducts(context.Background(), domain.Filter{
			Page:     1,
			PageSize: 2,
		})
		if err != nil {
			s.T().Fatal(err)
		}

		s.Suite.Assert().Equal(s.products, products)
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
		s.Suite.Assert().Equal([]*domain.Product{s.products[0]}, products)
		s.Suite.Assert().Equal(int64(2), n)

		n, products, err = s.db.GetProducts(context.Background(), domain.Filter{
			Page:     2,
			PageSize: 1,
		})
		if err != nil {
			s.T().Fatal(err)
		}

		s.Suite.Assert().Equal([]*domain.Product{s.products[1]}, products)
		s.Suite.Assert().Equal(int64(2), n)
	})
}

func (s *DatabaseTestSuite) SetupSuite() {
	s.setupContainer()
	s.setupAdapter()
	s.setupProducts()
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

	s.products = domainProducts(products)
}

func (s *DatabaseTestSuite) setupAdapter() {
	ctx := context.Background()

	db, err := NewDB(s.dsn)
	if err != nil {
		s.T().Fatalf("create new postgres connection pool: %s", err)
	}

	postgresDB := NewProductAdapter(db)

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
