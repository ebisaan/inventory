package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	inventoryv1 "github.com/ebisaan/proto/golang/ebisaan/inventory/v1beta1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ebisaan/inventory/config"
	"github.com/ebisaan/inventory/internal/application/port"
)

var _ inventoryv1.InventoryServiceServer = (*Adapter)(nil)

type Adapter struct {
	Done   chan struct{}
	server *grpc.Server
	app    port.API
	cfg    Config
	wg     sync.WaitGroup
	inventoryv1.UnimplementedInventoryServiceServer
}

type Config struct {
	Port int
	Env  string
}

func NewAdapter(api port.API, cfg Config) *Adapter {
	return &Adapter{
		app:  api,
		cfg:  cfg,
		Done: make(chan struct{}),
	}
}

func (a *Adapter) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on %d: %w", a.cfg.Port, err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(),
		),
		grpc.StreamInterceptor(
			recovery.StreamServerInterceptor(),
		),
	}

	srv := grpc.NewServer(opts...)
	inventoryv1.RegisterInventoryServiceServer(srv, a)
	if a.cfg.Env == config.DevEnv {
		reflection.Register(srv)
	}
	a.server = srv

	shutdownCh := make(chan struct{})
	go a.gracefulShutdown(shutdownCh)

	zap.L().Info(fmt.Sprintf("Starting gRPC server on port %d ...", a.cfg.Port))
	err = srv.Serve(l)
	if err != nil {
		return err
	}

	select {
	case <-time.After(10 * time.Second):
	case <-shutdownCh:
	}

	zap.L().Info("Stopped gRPC server")
	return nil
}

func (a *Adapter) gracefulShutdown(shutdown chan<- struct{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	zap.L().Info(fmt.Sprintf("Received signal %s", s))

	zap.L().Info("Shutdowning...")
	a.server.GracefulStop()

	zap.L().Info("Waiting for background tasks...")
	close(a.Done)
	a.wg.Wait()
	zap.L().Info("Background tasks completed")

	close(shutdown)
}

func (a *Adapter) Background(fn func()) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				zap.L().Error(fmt.Sprintf("Recovered from: %s", err))
			}
		}()

		fn()
	}()
}
