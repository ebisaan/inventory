package main

import (
	"context"
	"flag"

	"go.uber.org/zap"

	"github.com/ebisaan/inventory/config"
	"github.com/ebisaan/inventory/internal/adapter/grpc"
	"github.com/ebisaan/inventory/internal/adapter/postgres"
	"github.com/ebisaan/inventory/internal/application/core/app"
	"github.com/ebisaan/inventory/internal/logger"
)

func main() {
	var cfg config.Config
	filePath := flag.String("config", "config.yaml", "Configuration file")
	lgr, err := logger.New(zap.InfoLevel)
	if err != nil {
		zap.L().Fatal("failed to create logger: " + err.Error())
	}
	zap.ReplaceGlobals(lgr)

	flag.Parse()
	err = cfg.ReadFrom(*filePath)
	if err != nil {
		zap.L().Fatal("failed to read config from file: " + err.Error())
	}

	db, err := postgres.NewAdapter(cfg.DB.DSN)
	if err != nil {
		zap.L().Fatal("failed to create postgres adapter" + err.Error())
	}
	err = db.AutoMigration(context.Background())
	if err != nil {
		zap.L().Fatal("failed to migrate" + err.Error())
	}

	app := app.NewApplication(db)
	grpc := grpc.NewAdapter(app, grpc.Config{
		Port: cfg.Port,
		Env:  cfg.Env,
	})

	err = grpc.Run()
	if err != nil {
		zap.L().Fatal("failed to run grpc server" + err.Error())
	}
}
