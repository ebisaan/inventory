package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	inventory "github.com/ebisaan/proto/golang/inventory/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ebisaan/inventory/config"
	"github.com/ebisaan/inventory/internal/application/port"
)

var _ inventory.InventoryServiceServer = (*Adapter)(nil)

type Adapter struct {
	Done     chan struct{}
	server   *grpc.Server
	app      port.API
	cfg      Config
	wg       sync.WaitGroup
	shutdown chan struct{}
	inventory.UnimplementedInventoryServiceServer
}

type Config struct {
	Port int
	Env  string
}

func NewAdapter(api port.API, cfg Config) *Adapter {
	return &Adapter{
		app:      api,
		cfg:      cfg,
		Done:     make(chan struct{}),
		shutdown: make(chan struct{}),
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
	inventory.RegisterInventoryServiceServer(srv, a)
	if a.cfg.Env == config.DevEnv {
		reflection.Register(srv)
	}
	a.server = srv

	zap.L().Info(fmt.Sprintf("starting order grpc server on port %d ...", a.cfg.Port))
	go a.GracefulShutdown()
	err = srv.Serve(l)
	if err != nil {
		return err
	}

	select {
	case <-time.After(10 * time.Second):
	case <-a.shutdown:
	}

	return nil
}

func (a *Adapter) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	zap.L().Info(fmt.Sprintf("receive signal %s, stopping payment grpc server...", s))
	a.server.GracefulStop()
	close(a.Done)
	a.wg.Wait()
	close(a.shutdown)
}

func (a *Adapter) Background(fn func()) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				zap.L().Error(fmt.Sprintf("recover: %s", err))
			}
		}()

		fn()
	}()
}
