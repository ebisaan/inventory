package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/ebisaan/inventory/config"
	"github.com/ebisaan/inventory/internal/adapter/grpc"
	"github.com/ebisaan/inventory/internal/adapter/postgres"
	"github.com/ebisaan/inventory/internal/application/core/api"
	"github.com/ebisaan/inventory/internal/logger"
)

const configFile = "config.yaml"

func main() {
	var cfg config.Config

	lgr, err := logger.New(zap.InfoLevel)
	if err != nil {
		zap.L().Fatal("Failed to create logger: " + err.Error())
	}
	zap.ReplaceGlobals(lgr)

	err = cfg.ReadFrom(configFile)
	if err != nil {
		zap.L().Fatal("Failed to read config from file: " + err.Error())
	}

	parseFromFlags(&cfg)

	db, err := postgres.NewAdapter(cfg.DB.DSN, postgres.Config{
		MaxOpenConns: cfg.DB.MaxOpenConns,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxIdleTime:  cfg.DB.MaxIdleTime,
	})
	if err != nil {
		zap.L().Fatal("Failed to create postgres adapter" + err.Error())
	}

	app, err := api.NewApplication(db)
	if err != nil {
		zap.L().Fatal("Failed to create application adapter" + err.Error())
	}

	grpc := grpc.NewAdapter(app, grpc.Config{
		Port: cfg.Port,
		Env:  cfg.Env,
	})

	err = grpc.Run()
	if err != nil {
		zap.L().Fatal("Failed to run grpc server" + err.Error())
	}
}

func parseFromFlags(cfg *config.Config) {
	flag.Func("port", "API server's port", func(s string) error {
		var err error
		cfg.Port, err = strconv.Atoi(s)
		return err
	})
	flag.Func("dsn", "Data source name", func(s string) error {
		cfg.DB.DSN = s

		return nil
	})

	flag.Func("env", "Environment", func(s string) error {
		cfg.Env = s
		return nil
	})

	flag.Func("db-max-open-conns", "Max database open connections", func(s string) error {
		num, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid db-max-open-conns: %w", err)
		}
		cfg.DB.MaxOpenConns = num

		return nil
	})

	flag.Func("db-max-idle-conns", "Max database idle connections", func(s string) error {
		num, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid db-max-idle-conns: %w", err)
		}
		cfg.DB.MaxIdleConns = num

		return nil
	})

	flag.Func("db-max-idle-time", "Max database connection idle time", func(s string) error {
		dur, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid db-max-idle-time: %w", err)
		}
		cfg.DB.MaxIdleTime = dur

		return nil
	})

	flag.Parse()
}
