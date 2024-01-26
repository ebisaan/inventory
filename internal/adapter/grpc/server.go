package grpc

import (
	"fmt"
	"net"

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
	app port.API
	cfg Config
	inventory.UnimplementedInventoryServiceServer
}

type Config struct {
	Port int
	Env  string
}

func NewAdapter(api port.API, cfg Config) *Adapter {
	return &Adapter{
		app: api,
		cfg: cfg,
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

	zap.L().Info(fmt.Sprintf("starting order grpc server on port %d ...", a.cfg.Port))
	err = srv.Serve(l)
	if err != nil {
		return err
	}

	return nil
}
